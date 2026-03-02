#!/usr/bin/env node
/**
 * Import classic Sudoku puzzles from Cracking the Cryptic.
 *
 * For each puzzle:
 *   1. Fetch crackingthecryptic.com/sudoku?id=N to find the embedded
 *      SudokuPad token or tinyurl/short-link redirect.
 *   2. Follow redirects to resolve the final SudokuPad URL/token.
 *   3. Fetch sudokupad.app/api/puzzle/{token} — response may be either
 *      the old JS-object format  {ce:[[...]]}
 *      or the newer SCL format   "scl" + lz-string-base64-compressed JSON.
 *   4. Decode the grid; build an 81-char givens string.
 *   5. Skip non-standard puzzles (non-9x9, non-standard regions, extra
 *      constraints with numeric values such as killer cages).
 *   6. Emit JSON entries compatible with internal/curated/curated.json.
 */

const https = require('https');
const http  = require('http');

// ── lz-string decompressFromBase64 (verbatim from pieroxy/lz-string) ─────────
const _lzKeyB64 = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=';

function decompressFromBase64(input) {
  if (!input) return '';
  return _lzDecompress(input.length, 32, idx => _lzKeyB64.indexOf(input[idx]));
}

function _lzDecompress(length, resetValue, getNextValue) {
  var dictionary = [], next, enlargeIn = 4, dictSize = 4, numBits = 3,
      entry = '', result = [], i, w, bits, resb, maxpower, power, c,
      data = { val: getNextValue(0), position: resetValue, index: 1 };

  for (i = 0; i < 3; i++) dictionary[i] = i;

  bits = 0; maxpower = Math.pow(2, 2); power = 1;
  while (power !== maxpower) {
    resb = data.val & data.position;
    data.position >>= 1;
    if (data.position === 0) { data.position = resetValue; data.val = getNextValue(data.index++); }
    bits |= (resb > 0 ? 1 : 0) * power;
    power <<= 1;
  }

  switch (next = bits) {
    case 0:
      bits = 0; maxpower = Math.pow(2, 8); power = 1;
      while (power !== maxpower) {
        resb = data.val & data.position; data.position >>= 1;
        if (data.position === 0) { data.position = resetValue; data.val = getNextValue(data.index++); }
        bits |= (resb > 0 ? 1 : 0) * power; power <<= 1;
      }
      c = String.fromCharCode(bits); break;
    case 1:
      bits = 0; maxpower = Math.pow(2, 16); power = 1;
      while (power !== maxpower) {
        resb = data.val & data.position; data.position >>= 1;
        if (data.position === 0) { data.position = resetValue; data.val = getNextValue(data.index++); }
        bits |= (resb > 0 ? 1 : 0) * power; power <<= 1;
      }
      c = String.fromCharCode(bits); break;
    case 2:
      return '';
  }

  dictionary[3] = c; w = c; result.push(c);

  while (true) {
    if (data.index > length) return '';

    bits = 0; maxpower = Math.pow(2, numBits); power = 1;
    while (power !== maxpower) {
      resb = data.val & data.position; data.position >>= 1;
      if (data.position === 0) { data.position = resetValue; data.val = getNextValue(data.index++); }
      bits |= (resb > 0 ? 1 : 0) * power; power <<= 1;
    }

    switch (c = bits) {
      case 0:
        bits = 0; maxpower = Math.pow(2, 8); power = 1;
        while (power !== maxpower) {
          resb = data.val & data.position; data.position >>= 1;
          if (data.position === 0) { data.position = resetValue; data.val = getNextValue(data.index++); }
          bits |= (resb > 0 ? 1 : 0) * power; power <<= 1;
        }
        dictionary[dictSize++] = String.fromCharCode(bits); c = dictSize - 1; enlargeIn--; break;
      case 1:
        bits = 0; maxpower = Math.pow(2, 16); power = 1;
        while (power !== maxpower) {
          resb = data.val & data.position; data.position >>= 1;
          if (data.position === 0) { data.position = resetValue; data.val = getNextValue(data.index++); }
          bits |= (resb > 0 ? 1 : 0) * power; power <<= 1;
        }
        dictionary[dictSize++] = String.fromCharCode(bits); c = dictSize - 1; enlargeIn--; break;
      case 2:
        return result.join('');
    }

    if (enlargeIn === 0) { enlargeIn = Math.pow(2, numBits); numBits++; }

    if (dictionary[c]) { entry = dictionary[c]; }
    else { if (c === dictSize) { entry = w + w.charAt(0); } else { return null; } }
    result.push(entry);
    dictionary[dictSize++] = w + entry.charAt(0);
    enlargeIn--;
    if (enlargeIn === 0) { enlargeIn = Math.pow(2, numBits); numBits++; }
    w = entry;
  }
}

// ── HTTP helpers ─────────────────────────────────────────────────────────────
function fetchUrl(url, maxRedirects = 8) {
  return new Promise((resolve, reject) => {
    const attempt = (url, left) => {
      const mod = url.startsWith('https') ? https : http;
      const req = mod.get(url, {
        headers: { 'User-Agent': 'Mozilla/5.0 (compatible; sudoku-tui/1.0)' },
      }, res => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          res.resume();
          if (left <= 0) return reject(new Error('Too many redirects'));
          const next = new URL(res.headers.location, url).href;
          setTimeout(() => attempt(next, left - 1), 250);
          return;
        }
        if (res.statusCode !== 200) {
          res.resume();
          return reject(new Error(`HTTP ${res.statusCode}`));
        }
        let body = '';
        res.on('data', d => body += d);
        res.on('end', () => resolve({ url, body }));
      });
      req.on('error', reject);
      req.setTimeout(15000, () => { req.destroy(); reject(new Error('Timeout')); });
    };
    attempt(url, maxRedirects);
  });
}

const sleep = ms => new Promise(r => setTimeout(r, ms));

// ── Puzzle data decoding ──────────────────────────────────────────────────────

/**
 * Decode a SudokuPad API response.
 * Returns { format: 'jsobj'|'json', data: string|object } or null.
 */
function decodeApiResponse(raw) {
  raw = raw.trim();
  if (raw.startsWith('scl')) {
    const decompressed = decompressFromBase64(raw.slice(3));
    if (!decompressed) return null;
    // Decompressed may be JSON or JS-object-literal
    if (decompressed.startsWith('{') && decompressed.includes('"cells"')) {
      try { return { format: 'json', data: JSON.parse(decompressed) }; }
      catch { return { format: 'jsobj', data: decompressed }; }
    }
    return { format: 'jsobj', data: decompressed };
  }
  if (raw.startsWith('{')) return { format: 'jsobj', data: raw };
  return null;
}

// ── Grid extraction ───────────────────────────────────────────────────────────

/** Parse the JSON format: data.cells is 9x9 array of {value?: N} objects. */
function parseJsonCells(cells) {
  if (!Array.isArray(cells) || cells.length !== 9) return null;
  const grid = [];
  for (const row of cells) {
    if (!Array.isArray(row) || row.length !== 9) return null;
    grid.push(row.map(cell => (cell && typeof cell.value === 'number') ? cell.value : 0));
  }
  return grid;
}

/** Parse the JS-object-literal format: ce:[[...]] with {v:N} entries. */
function parseJsobjCE(jsObj) {
  const ceStart = jsObj.search(/\bce\s*:/);
  if (ceStart === -1) return null;

  const bracketStart = jsObj.indexOf('[', ceStart);
  if (bracketStart === -1) return null;

  let depth = 0, end = -1;
  for (let i = bracketStart; i < jsObj.length; i++) {
    if (jsObj[i] === '[') depth++;
    else if (jsObj[i] === ']') { depth--; if (depth === 0) { end = i + 1; break; } }
  }
  if (end === -1) return null;

  const ceStr = jsObj.slice(bracketStart, end);
  const rows = [];
  let rd = 0, rs = -1;
  for (let i = 1; i < ceStr.length - 1; i++) {
    if (ceStr[i] === '[' && rd === 0)  { rd = 1; rs = i + 1; }
    else if (ceStr[i] === '[')          rd++;
    else if (ceStr[i] === ']' && rd === 1) { rd = 0; rows.push(ceStr.slice(rs, i)); }
    else if (ceStr[i] === ']')          rd--;
  }
  if (rows.length !== 9) return null;

  return rows.map(rowStr => {
    const elements = [];
    let cur = '', bd = 0;
    for (const ch of rowStr) {
      if (ch === '{')          { bd++; cur += ch; }
      else if (ch === '}')     { bd--; cur += ch; }
      else if (ch === ',' && bd === 0) { elements.push(cur.trim()); cur = ''; }
      else cur += ch;
    }
    elements.push(cur.trim());
    const row = elements.map(e => { const m = e.match(/\{v\s*:\s*(\d+)\}/); return m ? parseInt(m[1]) : 0; });
    return row.length === 9 ? row : null;
  });
}

// ── Region validation ─────────────────────────────────────────────────────────
const STANDARD_REGIONS_SET = (() => {
  const regs = [];
  for (let br = 0; br < 3; br++)
    for (let bc = 0; bc < 3; bc++) {
      const cells = [];
      for (let r = 0; r < 3; r++)
        for (let c = 0; c < 3; c++)
          cells.push(br * 3 + r, bc * 3 + c);   // interleaved: [r, c, r, c, ...]
      regs.push(cells);
    }
  return regs;
})();

function canonicalRegion(pairs) {
  // pairs is array of [row,col] or interleaved [r,c,r,c,...]
  let cells;
  if (Array.isArray(pairs[0])) {
    cells = pairs.map(([r, c]) => r * 9 + c);
  } else {
    cells = [];
    for (let i = 0; i < pairs.length; i += 2) cells.push(pairs[i] * 9 + pairs[i + 1]);
  }
  return cells.sort((a, b) => a - b).join(',');
}

const STANDARD_REGION_KEYS = new Set(
  STANDARD_REGIONS_SET.map((_, bi) => {
    const cells = [];
    const br = Math.floor(bi / 3), bc = bi % 3;
    for (let r = 0; r < 3; r++)
      for (let c = 0; c < 3; c++)
        cells.push((br * 3 + r) * 9 + (bc * 3 + c));
    return cells.sort((a, b) => a - b).join(',');
  })
);

function isStandardRegions(decoded) {
  let regionsList;

  if (decoded.format === 'json') {
    const data = decoded.data;
    if (!data.regions) return true; // no regions = standard
    regionsList = data.regions;
  } else {
    const jsObj = decoded.data;
    const reStart = jsObj.search(/\bre\s*:/);
    if (reStart === -1) return true;
    const bracketStart = jsObj.indexOf('[', reStart);
    if (bracketStart === -1) return true;
    let depth = 0, end = -1;
    for (let i = bracketStart; i < jsObj.length; i++) {
      if (jsObj[i] === '[') depth++;
      else if (jsObj[i] === ']') { depth--; if (depth === 0) { end = i + 1; break; } }
    }
    if (end === -1) return true;
    try { regionsList = JSON.parse(jsObj.slice(bracketStart, end)); }
    catch { return true; } // can't parse, assume standard
  }

  if (!Array.isArray(regionsList) || regionsList.length !== 9) return false;
  const keys = new Set(regionsList.map(canonicalRegion));
  return keys.size === 9 && [...keys].every(k => STANDARD_REGION_KEYS.has(k));
}

/**
 * Returns true if the puzzle has extra constraints that make it non-classic.
 * We check for:
 *  - Killer cages (cages with a numeric value/sum)
 *  - Thermo, arrow, sandwich, diagonal fields
 */
function hasExtraConstraints(decoded) {
  if (decoded.format === 'json') {
    const data = decoded.data;
    // Killer cages have numeric sum values
    if (Array.isArray(data.cages)) {
      for (const cage of data.cages) {
        if (typeof cage.value === 'number') return true;
        if (typeof cage.value === 'string' && /^\d+$/.test(cage.value)) return true;
      }
    }
    // Lines = thermo/arrow/etc.
    if (data.lines && data.lines.length > 0) return true;
    // Other known constraint fields
    const constraintKeys = ['thermometers', 'arrows', 'sandwichSums', 'diagonals',
                            'xSum', 'killercages', 'extraRegions'];
    if (constraintKeys.some(k => data[k] && data[k].length > 0)) return true;
    return false;
  } else {
    // JS object literal: check for non-metadata constraint patterns
    const jsObj = decoded.data;
    // Known constraint field names in JS object notation
    return /\b(tc|kc|ac|sc|di|lines|thermometers|arrows|sandwichSums)\s*:\s*\[/.test(jsObj);
  }
}

// ── Token extraction from CTC page ───────────────────────────────────────────

async function getSudokuPadInfo(ctcId) {
  let html;
  try {
    const r = await fetchUrl(`https://crackingthecryptic.com/sudoku?id=${ctcId}`);
    html = r.body;
  } catch (e) {
    process.stderr.write(`  Fetch error: ${e.message}\n`);
    return null;
  }

  // Direct SudokuPad embed: https://sudokupad.app/TOKEN
  const directMatch = html.match(/sudokupad\.app\/([A-Za-z0-9_\-+=/.]+)/);
  if (directMatch) {
    const val = directMatch[1].replace(/\/$/, '');
    return { type: val.startsWith('scl') ? 'scl' : 'token', value: val };
  }

  // Short-link / tinyurl redirect
  const shortMatch = html.match(/https?:\/\/(tinyurl\.com|bit\.ly)\/([A-Za-z0-9_-]+)/);
  if (shortMatch) {
    try {
      const r2 = await fetchUrl(shortMatch[0]);
      const final = r2.url;
      process.stderr.write(`  redirect → ${final.slice(0, 80)}\n`);
      const m2 = final.match(/sudokupad\.app\/([A-Za-z0-9_\-+=/.]+)/);
      if (m2) {
        const val = m2[1].replace(/\/$/, '');
        return { type: val.startsWith('scl') ? 'scl' : 'token', value: val };
      }
    } catch (e) {
      process.stderr.write(`  short-link error: ${e.message}\n`);
    }
  }

  return null;
}

async function getPuzzleDecoded(info) {
  if (info.type === 'scl') {
    // The SCL value IS the encoded puzzle
    const decompressed = decompressFromBase64(info.value.slice(3));
    if (!decompressed) return null;
    if (decompressed.startsWith('{') && decompressed.includes('"cells"')) {
      try { return { format: 'json', data: JSON.parse(decompressed) }; }
      catch { return { format: 'jsobj', data: decompressed }; }
    }
    return { format: 'jsobj', data: decompressed };
  }

  // token: fetch from API
  try {
    const r = await fetchUrl(`https://sudokupad.app/api/puzzle/${info.value}`);
    return decodeApiResponse(r.body);
  } catch (e) {
    process.stderr.write(`  API fetch error: ${e.message}\n`);
    return null;
  }
}

// ── Entry construction ────────────────────────────────────────────────────────
function difficultyFromCount(n) {
  if (n >= 36) return 'easy';
  if (n >= 28) return 'medium';
  if (n >= 22) return 'hard';
  return 'very hard';
}

function makeId(ctcId, title) {
  const safe = title.toLowerCase().replace(/[^a-z0-9]+/g, '_').replace(/^_|_$/g, '').slice(0, 30);
  return `ctc_${ctcId}_${safe}`;
}

// ── Puzzle list ───────────────────────────────────────────────────────────────
const PUZZLES = [
  [2302, "The Loneliest Number", "Snyder"],
  [2250, "Crest of the Gorons", "Wei-Hwa Huang"],
  [2246, "Din's Pearl", "Wei-Hwa Huang"],
  [2247, "Crest of the Kokiri", "Wei-Hwa Huang"],
  [2248, "Forest Medallion", "Wei-Hwa Huang"],
  [485, "White Marlin", "Shye"],
  [436, "Arbitrary Code Execution", "Jovial"],
  [435, "Can't Teach An Old Dog", "Jovial"],
  [484, "Antidote", "Shye"],
  [483, "Smile", "Shye"],
  [467, "pipeline", "Mith"],
  [1581, "Long Distance Relationship", "RiSa"],
  [1279, "Fancy Vase", "Qinlux"],
  [481, "Kingda Ka", "Shye"],
  [480, "Binary Fission", "Shye"],
  [668, "Computer Solver Freaks Out", "Feadoor"],
  [1278, "noXing", "Qinlux"],
  [478, "Patto Patto", "Shye"],
  [206, "Cobra Roll", "Jovial"],
  [1301, "Steering Wheel", "SudokuExplorer"],
  [667, "The XY Ring", "Feadoor"],
  [421, "Hanabi", "Shye"],
  [419, "Valtari", "Shye"],
  [1626, "Give Me an X, Give Me a V", "aPete"],
  [683, "Learning to Solve Hard Sudoku", "Doulani"],
  [200, "Classic Sudoku", "Jovial"],
  [799, "Gridlocked", "Snyder"],
  [464, "22 Clues 7 Rows", "Mith"],
  [2447, "Skirmish in the Horsehead Nebula", "Øyvind"],
  [664, "The Expert's Sudoku Technique", "Feadoor"],
  [418, "Monte", "Shye"],
  [662, "Teaching the X-wing", "Feadoor"],
  [661, "Human Logic Trumps Computer", "Feadoor"],
  [459, "Tatooine Sunset", "Mith"],
  [660, "Stripe Sudoku", "Feadoor"],
  [995, "How do Grandmasters Solve so Fast", "Collyer"],
  [1326, "4x4 Magic Square", "Cam"],
  [559, "SVS(266) Tic Tac Toe", "Richard"],
  [658, "Unbelievable Masterpiece", "Feadoor"],
  [1324, "Goldilocks Sudoku", "Cam"],
  [657, "Small but Beautiful", "Feadoor"],
  [655, "How to find X-wings", "Feadoor"],
  [1153, "Bonus Linked No. 4", "Kumar"],
  [1152, "Bonus Linked No. 3", "Kumar"],
  [1151, "Bonus Linked No. 2", "Kumar"],
  [1150, "Bonus Linked No. 1", "Kumar"],
  [652, "Hidden Pairs", "Feadoor"],
  [1526, "Beyond Diabolical", "DerekNeal"],
  [1149, "8 Minute Clone", "Kumar"],
  [501, "Pi Sudoku", "Aad"],
  [651, "Hardest Sudoku", "Feadoor"],
  [497, "The Roaring Four-T's", "Aad"],
  [1080, "In Memoriam", "Holaysan"],
  [496, "T Sudoku", "Aad"],
  [1525, "The Stinger", "DerekNeal"],
  [495, "<=5 Sudoku", "Aad"],
  [718, "Less than 15 Sudoku", "Bastien"],
  [795, "Three Ring Circus", "Snyder"],
  [649, "Hard in 10 Minutes", "Feadoor"],
  [1524, "Celebrating 50k Subs", "DerekNeal"],
  [714, "Spot the Initial Trick", "Bastien"],
  [1147, "The Schizophrenic Sudoku", "Kumar"],
  [712, "Stunning Sudoku from France", "Bastien"],
  [494, "A Dutch Masterpiece", "Aad"],
  [585, "Classic", "Richard"],
  [648, "Jellyfish", "Feadoor"],
  [646, "Extreme Sudoku made Easy", "Feadoor"],
  [645, "How Many X-Wings", "Feadoor"],
  [1523, "Diabolical Excalibur!", "DerekNeal"],
  [644, "Diabolical Sudoku", "Feadoor"],
  [1522, "Diabolical Simple Trick", "DerekNeal"],
  [1521, "Hidden Triple", "DerekNeal"],
  [2245, "Antidiagonal", "Wei-Hwa Huang"],
  [643, "Find that Shape", "Feadoor"],
  [914, "Top Heavy", "Yürekli"],
  [56, "Diagonal Sudoku", "Collyer"],
  [790, "Diagonals", "Snyder"],
];

// ── Main ──────────────────────────────────────────────────────────────────────
async function main() {
  const results = [];
  const skipped = [];

  for (const [ctcId, title, author] of PUZZLES) {
    process.stderr.write(`\n[${ctcId}] ${title} (${author})\n`);

    const info = await getSudokuPadInfo(ctcId);
    if (!info) {
      process.stderr.write(`  SKIP: no SudokuPad info\n`);
      skipped.push([ctcId, title, 'no SudokuPad info']); await sleep(400); continue;
    }
    process.stderr.write(`  type=${info.type} val=${info.value.slice(0, 25)}...\n`);
    await sleep(300);

    const decoded = await getPuzzleDecoded(info);
    if (!decoded) {
      process.stderr.write(`  SKIP: could not decode puzzle data\n`);
      skipped.push([ctcId, title, 'decode failed']); await sleep(400); continue;
    }

    if (!isStandardRegions(decoded)) {
      process.stderr.write(`  SKIP: non-standard regions\n`);
      skipped.push([ctcId, title, 'non-standard regions']); await sleep(200); continue;
    }

    if (hasExtraConstraints(decoded)) {
      process.stderr.write(`  SKIP: extra constraints\n`);
      skipped.push([ctcId, title, 'extra constraints']); await sleep(200); continue;
    }

    let grid;
    if (decoded.format === 'json') {
      grid = parseJsonCells(decoded.data.cells);
    } else {
      grid = parseJsobjCE(decoded.data);
    }

    if (!grid || grid.some(row => !row || row.length !== 9)) {
      process.stderr.write(`  SKIP: grid parse failed (format=${decoded.format})\n`);
      skipped.push([ctcId, title, 'grid parse failed']); await sleep(200); continue;
    }

    const givens = grid.flat().join('');
    if (givens.length !== 81) {
      process.stderr.write(`  SKIP: givens length ${givens.length}\n`);
      skipped.push([ctcId, title, `bad givens length`]); continue;
    }

    const nClues = givens.split('').filter(c => c !== '0').length;
    const diff = difficultyFromCount(nClues);
    process.stderr.write(`  OK: format=${decoded.format} clues=${nClues} → ${diff}\n`);

    results.push({
      id: makeId(ctcId, title),
      name: title,
      difficulty: diff,
      author: author,
      reference: `https://crackingthecryptic.com/sudoku?id=${ctcId}`,
      description: `Classic sudoku by ${author}, featured on Cracking the Cryptic. ${nClues} given digits.`,
      givens: givens,
    });

    await sleep(300);
  }

  process.stderr.write(`\nImported: ${results.length}, Skipped: ${skipped.length}\n`);
  if (skipped.length > 0) {
    process.stderr.write('Skipped:\n');
    skipped.forEach(([id, name, reason]) =>
      process.stderr.write(`  [${id}] ${name}: ${reason}\n`));
  }

  process.stdout.write(JSON.stringify(results, null, 2) + '\n');
}

main().catch(e => { process.stderr.write(`Fatal: ${e}\n`); process.exit(1); });
