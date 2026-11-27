# The Python helper of GAN's python(...) function. It is built into gan.exe (see python.go).
#
# GAN starts it once, then sends one request per line and reads one reply per line (JSON).
#
# It can call ANY function in Python's modules, with the same power as any Python program.
import sys

# Python normally looks for modules in the current folder first. Take that out BEFORE the
# first import, so a stray file such as json.py or numpy.py in the folder is never run by
# mistake.
sys.path = [p for p in sys.path if p not in ("", ".")]

import importlib
import json
import math
import numbers
import os

# Keep the two pipes to GAN for ourselves. After this, anything that Python code prints
# appears on the screen, and code that tries to read input gets nothing, so neither can
# break the conversation with GAN.
from_gan = os.fdopen(os.dup(0), "r", encoding="utf-8", newline="\n")
to_gan = os.fdopen(os.dup(1), "w", encoding="utf-8", newline="\n")
nothing = os.open(os.devnull, os.O_RDONLY)
os.dup2(nothing, 0)
os.dup2(2, 1)
sys.stdin = open(os.devnull, "r")
sys.stdout = sys.stderr

MAX_DEPTH = 100             # lists inside lists
MAX_ELEMENTS = 10_000_000   # the same limit GAN has for tensors


class Problem(Exception):
    """A mistake to report to the GAN program, already in plain words."""


def load_numpy():
    try:
        return importlib.import_module("numpy")
    except ImportError:
        raise Problem("sending a tensor to Python needs NumPy, and Python could not load it")


def from_value(v):
    """A GAN value (as sent by GAN) -> a Python value."""
    kind = v["type"]
    if kind == "number":
        x = float(v["value"])
        # whole numbers arrive as Python integers, because many functions insist on that
        return int(x) if x.is_integer() and abs(x) <= 2 ** 53 else x
    if kind in ("text", "bool"):
        return v["value"]
    if kind == "list":
        return [from_value(item) for item in v["items"]]
    numpy = load_numpy()
    return numpy.array([float(x) for x in v["data"]], dtype=float).reshape(v["shape"])


def number(x):
    try:
        return {"type": "number", "value": repr(float(x))}
    except OverflowError:
        raise Problem("Python returned a number that is too big for GAN")


def to_value(x, depth=0):
    """A Python result -> a GAN value. None means "no value"."""
    if depth > MAX_DEPTH:
        raise Problem("Python returned lists nested too deeply for GAN")
    numpy = sys.modules.get("numpy")
    if numpy is not None:
        if isinstance(x, numpy.ndarray):
            if x.dtype.kind not in "biuf":
                raise Problem("Python returned a NumPy array of " + str(x.dtype) + ", but GAN tensors only hold numbers")
            if x.ndim == 0:
                return to_value(x.item(), depth)
            if x.size == 0:
                raise Problem("Python returned an empty array, and a GAN tensor cannot be empty")
            if x.size > MAX_ELEMENTS or x.ndim > 32:
                raise Problem("Python returned an array that is too big for a GAN tensor")
            return {"type": "tensor", "shape": list(x.shape), "data": [repr(float(n)) for n in x.ravel().tolist()]}
        if isinstance(x, numpy.generic):
            x = x.item()
    if isinstance(x, bool):
        return {"type": "bool", "value": x}
    if isinstance(x, numbers.Real):
        return number(x)
    if isinstance(x, str):
        return {"type": "text", "value": x}
    if isinstance(x, (list, tuple)):
        items = [to_value(item, depth + 1) for item in x]
        if any(item is None for item in items):
            raise Problem("Python returned a list with None in it, which GAN cannot use yet")
        return {"type": "list", "items": items}
    if x is None:
        return None
    raise Problem("Python returned a " + type(x).__name__ + ", which GAN cannot use yet")


def find_function(name):
    """ "numpy.linalg.inv" -> import numpy.linalg and fetch inv. """
    parts = name.split(".")
    for split in range(len(parts) - 1, 0, -1):
        module_name = ".".join(parts[:split])
        try:
            found = importlib.import_module(module_name)
        except ModuleNotFoundError as err:
            # "not found" for this name (or a parent of it): try a shorter module name
            if err.name and (module_name == err.name or module_name.startswith(err.name + ".")):
                continue
            raise  # the module exists, but something IT needs is missing
        for part in parts[split:]:
            if not hasattr(found, part):
                raise Problem('Python has no function "' + name + '"')
            found = getattr(found, part)
        if not callable(found):
            raise Problem('"' + name + '" in Python is not a function')
        return found
    raise Problem('Python has no module "' + parts[0] + '"; is it installed?')


def handle(request):
    name = request["function"]
    try:
        function = find_function(name)
        result = function(*[from_value(a) for a in request["args"]])
        return {"ok": True, "value": to_value(result)}
    except Problem as err:
        return {"ok": False, "error": str(err)}
    except KeyboardInterrupt:
        raise
    except BaseException as err:  # includes sys.exit(), which must not end the helper silently
        message = " ".join(str(err).split())  # one line
        return {"ok": False, "error": "Python error in " + name + ": " + type(err).__name__ + (": " + message if message else "")}


def send(reply):
    to_gan.write(json.dumps(reply, allow_nan=False) + "\n")
    to_gan.flush()


send({"ready": True})
for line in from_gan:
    reply = handle(json.loads(line))
    try:
        sys.stderr.flush()  # so that what Python printed appears before GAN carries on
    except Exception:
        pass
    send(reply)
