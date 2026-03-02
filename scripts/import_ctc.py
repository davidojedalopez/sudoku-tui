#!/usr/bin/env python3
"""
Import classic Sudoku puzzles from Cracking the Cryptic.

For each puzzle:
  1. Fetch crackingthecryptic.com/sudoku?id=N to find the embedded SudokuPad token.
  2. Fetch sudokupad.app/api/puzzle/{token} to get the puzzle data.
  3. Parse the 'ce' (cell entries) array to build an 81-char givens string.
  4. Skip puzzles that aren't standard 9x9 classics.
  5. Write JSON entries compatible with internal/curated/curated.json.
"""

import re
import time
import json
import sys
import urllib.request
import urllib.error

# All 77 CTC classic puzzles: (id, title, author)
PUZZLES = [
    (2302, "The Loneliest Number", "Snyder"),
    (2250, "Crest of the Gorons", "Wei-Hwa Huang"),
    (2246, "Din's Pearl", "Wei-Hwa Huang"),
    (2247, "Crest of the Kokiri", "Wei-Hwa Huang"),
    (2248, "Forest Medallion", "Wei-Hwa Huang"),
    (485, "White Marlin", "Shye"),
    (436, "Arbitrary Code Execution", "Jovial"),
    (435, "Can't Teach An Old Dog", "Jovial"),
    (484, "Antidote", "Shye"),
    (483, "Smile", "Shye"),
    (467, "pipeline", "Mith"),
    (1581, "Long Distance Relationship", "RiSa"),
    (1279, "Fancy Vase", "Qinlux"),
    (481, "Kingda Ka", "Shye"),
    (480, "Binary Fission", "Shye"),
    (668, "Computer Solver Freaks Out", "Feadoor"),
    (1278, "noXing", "Qinlux"),
    (478, "Patto Patto", "Shye"),
    (206, "Cobra Roll", "Jovial"),
    (1301, "Steering Wheel", "SudokuExplorer"),
    (667, "The XY Ring", "Feadoor"),
    (421, "Hanabi", "Shye"),
    (419, "Valtari", "Shye"),
    (1626, "Give Me an X, Give Me a V", "aPete"),
    (683, "Learning to Solve Hard Sudoku", "Doulani"),
    (200, "Classic Sudoku", "Jovial"),
    (799, "Gridlocked", "Snyder"),
    (464, "22 Clues 7 Rows", "Mith"),
    (2447, "Skirmish in the Horsehead Nebula", "Øyvind"),
    (664, "The Expert's Sudoku Technique", "Feadoor"),
    (418, "Monte", "Shye"),
    (662, "Teaching the X-wing", "Feadoor"),
    (661, "Human Logic Trumps Computer", "Feadoor"),
    (459, "Tatooine Sunset", "Mith"),
    (660, "Stripe Sudoku", "Feadoor"),
    (995, "How do Grandmasters Solve so Fast", "Collyer"),
    (1326, "4x4 Magic Square", "Cam"),
    (559, "SVS(266) Tic Tac Toe", "Richard"),
    (658, "Unbelievable Masterpiece", "Feadoor"),
    (1324, "Goldilocks Sudoku", "Cam"),
    (657, "Small but Beautiful", "Feadoor"),
    (655, "How to find X-wings", "Feadoor"),
    (1153, "Bonus Linked No. 4", "Kumar"),
    (1152, "Bonus Linked No. 3", "Kumar"),
    (1151, "Bonus Linked No. 2", "Kumar"),
    (1150, "Bonus Linked No. 1", "Kumar"),
    (652, "Hidden Pairs", "Feadoor"),
    (1526, "Beyond Diabolical", "DerekNeal"),
    (1149, "8 Minute Clone", "Kumar"),
    (501, "Pi Sudoku", "Aad"),
    (651, "Hardest Sudoku", "Feadoor"),
    (497, "The Roaring Four-T's", "Aad"),
    (1080, "In Memoriam", "Holaysan"),
    (496, "T Sudoku", "Aad"),
    (1525, "The Stinger", "DerekNeal"),
    (495, "<=5 Sudoku", "Aad"),
    (718, "Less than 15 Sudoku", "Bastien"),
    (795, "Three Ring Circus", "Snyder"),
    (649, "Hard in 10 Minutes", "Feadoor"),
    (1524, "Celebrating 50k Subs", "DerekNeal"),
    (714, "Spot the Initial Trick", "Bastien"),
    (1147, "The Schizophrenic Sudoku", "Kumar"),
    (712, "Stunning Sudoku from France", "Bastien"),
    (494, "A Dutch Masterpiece", "Aad"),
    (585, "Classic", "Richard"),
    (648, "Jellyfish", "Feadoor"),
    (646, "Extreme Sudoku made Easy", "Feadoor"),
    (645, "How Many X-Wings", "Feadoor"),
    (1523, "Diabolical Excalibur!", "DerekNeal"),
    (644, "Diabolical Sudoku", "Feadoor"),
    (1522, "Diabolical Simple Trick", "DerekNeal"),
    (1521, "Hidden Triple", "DerekNeal"),
    (2245, "Antidiagonal", "Wei-Hwa Huang"),
    (643, "Find that Shape", "Feadoor"),
    (914, "Top Heavy", "Yürekli"),
    (56, "Diagonal Sudoku", "Collyer"),
    (790, "Diagonals", "Snyder"),
]

STANDARD_REGIONS = [
    [[0,0],[0,1],[0,2],[1,0],[1,1],[1,2],[2,0],[2,1],[2,2]],
    [[3,0],[3,1],[3,2],[4,0],[4,1],[4,2],[5,0],[5,1],[5,2]],
    [[6,0],[6,1],[6,2],[7,0],[7,1],[7,2],[8,0],[8,1],[8,2]],
    [[0,3],[0,4],[0,5],[1,3],[1,4],[1,5],[2,3],[2,4],[2,5]],
    [[3,3],[3,4],[3,5],[4,3],[4,4],[4,5],[5,3],[5,4],[5,5]],
    [[6,3],[6,4],[6,5],[7,3],[7,4],[7,5],[8,3],[8,4],[8,5]],
    [[0,6],[0,7],[0,8],[1,6],[1,7],[1,8],[2,6],[2,7],[2,8]],
    [[3,6],[3,7],[3,8],[4,6],[4,7],[4,8],[5,6],[5,7],[5,8]],
    [[6,6],[6,7],[6,8],[7,6],[7,7],[7,8],[8,6],[8,7],[8,8]],
]


def fetch(url, retries=3):
    for attempt in range(retries):
        try:
            req = urllib.request.Request(url, headers={
                "User-Agent": "Mozilla/5.0 (compatible; sudoku-tui-importer/1.0)"
            })
            with urllib.request.urlopen(req, timeout=15) as resp:
                return resp.read().decode("utf-8", errors="replace")
        except urllib.error.HTTPError as e:
            print(f"  HTTP {e.code} for {url}", file=sys.stderr)
            if e.code == 404:
                return None
            time.sleep(2)
        except Exception as e:
            print(f"  Error fetching {url}: {e}", file=sys.stderr)
            time.sleep(2)
    return None


def get_sudokupad_token(ctc_id):
    """Fetch CTC puzzle page and extract the embedded SudokuPad token."""
    url = f"https://crackingthecryptic.com/sudoku?id={ctc_id}"
    html = fetch(url)
    if not html:
        return None
    # Look for https://sudokupad.app/XXXXX link
    m = re.search(r'https://sudokupad\.app/([A-Za-z0-9_-]+)', html)
    if m:
        token = m.group(1)
        # Ignore paths that look like routes (contain slashes handled above,
        # or are known non-puzzle paths)
        if token not in ("sudoku", "api", "puzzle"):
            return token
    return None


def get_puzzle_data(token):
    """Fetch puzzle data from SudokuPad API."""
    url = f"https://sudokupad.app/api/puzzle/{token}"
    return fetch(url)


def parse_ce(raw):
    """
    Parse the 'ce' field from SudokuPad JS object notation.
    Returns a 9x9 list of ints (0 = empty).
    """
    # Extract the ce value: everything after 'ce:' up to the matching bracket
    m = re.search(r'\bce:\s*(\[)', raw)
    if not m:
        return None

    start = m.start(1)
    depth = 0
    end = start
    for i in range(start, len(raw)):
        if raw[i] == '[':
            depth += 1
        elif raw[i] == ']':
            depth -= 1
            if depth == 0:
                end = i + 1
                break

    ce_str = raw[start:end]

    # Extract each row array [...] from the outer [[...],[...],...]
    # We iterate over the top-level contents
    # Strip outer brackets
    inner = ce_str[1:-1]  # remove outer [ and ]

    rows = []
    depth = 0
    row_start = None
    for i, c in enumerate(inner):
        if c == '[' and depth == 0:
            depth = 1
            row_start = i + 1
        elif c == '[':
            depth += 1
        elif c == ']' and depth == 1:
            depth = 0
            rows.append(inner[row_start:i])
        elif c == ']':
            depth -= 1

    if len(rows) != 9:
        return None

    grid = []
    for row_str in rows:
        # Split by comma respecting nested braces
        elements = []
        cur = ""
        bdepth = 0
        for c in row_str:
            if c == '{':
                bdepth += 1
                cur += c
            elif c == '}':
                bdepth -= 1
                cur += c
            elif c == ',' and bdepth == 0:
                elements.append(cur.strip())
                cur = ""
            else:
                cur += c
        elements.append(cur.strip())

        row = []
        for elem in elements:
            vm = re.match(r'\{v:(\d+)\}', elem)
            row.append(int(vm.group(1)) if vm else 0)

        if len(row) != 9:
            return None
        grid.append(row)

    return grid


def parse_regions(raw):
    """
    Parse the 're' field and check if it matches standard 9x3 box regions.
    Returns True if standard, False otherwise.
    """
    m = re.search(r'\bre:\s*(\[)', raw)
    if not m:
        # No regions field; assume standard
        return True

    start = m.start(1)
    depth = 0
    end = start
    for i in range(start, len(raw)):
        if raw[i] == '[':
            depth += 1
        elif raw[i] == ']':
            depth -= 1
            if depth == 0:
                end = i + 1
                break

    re_str = raw[start:end]

    # Extract all [row,col] pairs grouped by region
    region_pattern = re.findall(r'\[(\d+),(\d+)\]', re_str)
    if len(region_pattern) != 81:
        return False

    # Build the 9 regions as sorted lists of (row, col)
    parsed_regions = []
    # Re-parse grouping by the outer brackets
    inner = re_str[1:-1]
    regions = []
    depth = 0
    reg_start = None
    for i, c in enumerate(inner):
        if c == '[' and depth == 0:
            depth = 1
            reg_start = i
        elif c == '[':
            depth += 1
        elif c == ']' and depth == 1:
            depth = 0
            regions.append(inner[reg_start:i+1])
        elif c == ']':
            depth -= 1

    if len(regions) != 9:
        return False

    parsed = []
    for reg_str in regions:
        cells = re.findall(r'\[(\d+),(\d+)\]', reg_str)
        parsed.append(sorted([(int(r), int(c)) for r, c in cells]))

    expected = [sorted([(r, c) for r, c in reg]) for reg in STANDARD_REGIONS]
    return sorted(parsed) == sorted(expected)


def grid_to_givens(grid):
    return "".join(str(v) for row in grid for v in row)


def count_givens(givens):
    return sum(1 for c in givens if c != '0')


def difficulty_from_count(n):
    if n >= 36:
        return "easy"
    elif n >= 28:
        return "medium"
    elif n >= 22:
        return "hard"
    else:
        return "very hard"


def make_id(ctc_id, title):
    safe = re.sub(r'[^a-z0-9]+', '_', title.lower()).strip('_')
    return f"ctc_{ctc_id}_{safe[:30]}"


def main():
    results = []
    skipped = []

    for ctc_id, title, author in PUZZLES:
        print(f"[{ctc_id}] {title} ({author})", file=sys.stderr)

        token = get_sudokupad_token(ctc_id)
        if not token:
            print(f"  SKIP: no SudokuPad token found", file=sys.stderr)
            skipped.append((ctc_id, title, "no token"))
            time.sleep(0.5)
            continue

        print(f"  token: {token}", file=sys.stderr)
        time.sleep(0.3)

        raw = get_puzzle_data(token)
        if not raw:
            print(f"  SKIP: could not fetch puzzle data", file=sys.stderr)
            skipped.append((ctc_id, title, "no puzzle data"))
            time.sleep(0.5)
            continue

        if not parse_regions(raw):
            print(f"  SKIP: non-standard regions (not a classic 9x9)", file=sys.stderr)
            skipped.append((ctc_id, title, "non-standard regions"))
            time.sleep(0.3)
            continue

        grid = parse_ce(raw)
        if not grid:
            print(f"  SKIP: could not parse ce field", file=sys.stderr)
            skipped.append((ctc_id, title, "parse error"))
            time.sleep(0.3)
            continue

        givens = grid_to_givens(grid)
        if len(givens) != 81:
            print(f"  SKIP: givens length {len(givens)} != 81", file=sys.stderr)
            skipped.append((ctc_id, title, f"bad givens len={len(givens)}"))
            time.sleep(0.3)
            continue

        n_clues = count_givens(givens)
        diff = difficulty_from_count(n_clues)
        print(f"  givens: {n_clues} clues → {diff}", file=sys.stderr)

        entry = {
            "id": make_id(ctc_id, title),
            "name": title,
            "difficulty": diff,
            "author": author,
            "reference": f"https://crackingthecryptic.com/sudoku?id={ctc_id}",
            "description": f"Classic sudoku by {author}, featured on Cracking the Cryptic. {n_clues} given digits.",
            "givens": givens,
        }
        results.append(entry)
        time.sleep(0.3)

    print(f"\nImported: {len(results)}, Skipped: {len(skipped)}", file=sys.stderr)
    if skipped:
        print("Skipped puzzles:", file=sys.stderr)
        for ctc_id, title, reason in skipped:
            print(f"  [{ctc_id}] {title}: {reason}", file=sys.stderr)

    print(json.dumps(results, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
