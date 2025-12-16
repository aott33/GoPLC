# GoPLC Logo Improvement Guidelines

This document captures the agreed-upon improvements to refine the GoPLC logo, with a focus on clarity, proportion, and usability across contexts.

---

## Goals

* Make the mascot **unambiguously a gopher** (especially for Go developers)
* Improve **overall proportions** and move the logo closer to a square
* Increase **readability at small sizes** (icons, favicons, CLI, README)
* Preserve the existing visual identity and tone

---

## 1. Clarify the Gopher Identity (High Priority)

The current head shape is close, but slightly ambiguous. Small, targeted cues will clearly communicate "gopher" without a redesign.

### Recommended Gopher Cues

Apply **2–3 of the following** (do not use all):

* **Eye**

  * Slightly larger and more circular
  * Optional small white highlight dot

* **Tooth**

  * Add a small white rectangular tooth at the front of the snout
  * This is one of the strongest rodent/gopher signals

* **Snout Shape**

  * Slightly round the snout edges
  * Reduce sharp angles at the mouth

* **Ear Hint (Optional)**

  * Very subtle bump or notch behind the eye

> Even a single tooth plus a slightly larger eye will immediately read as a gopher to Go developers.

---

## 2. Adjust Gopher Size & Proportion (High Priority)

The gopher currently dominates the vertical space.

### Size Adjustment

* Reduce gopher height by **15–25%**
* Keep width mostly unchanged
* Lower the gopher slightly so it feels seated on the PLC panel

### Target Vertical Ratio

* **Gopher:** ~45% of total height
* **PLC panel:** ~55% of total height

This will:

* Reduce top-heaviness
* Improve balance
* Move the logo closer to a square

---

## 3. Square-Friendly Composition

The logo should fit naturally inside a square for avatars, icons, and badges.

### Recommendations

* Ensure similar visual margins at the top (gopher) and bottom (buttons)
* If needed:

  * Slightly reduce the height of the bottom button row
  * Tighten vertical padding inside the PLC panel

Goal: minimal dead space when centered in a square.

---

## 4. Seat the Gopher on the PLC

To visually connect the mascot and the hardware concept:

* Flatten the bottom edge of the gopher slightly
* Align it clearly with the top edge of the PLC panel
* Avoid the appearance that the gopher is floating above the panel

This reinforces the concept:

> **Go (gopher) + PLC (hardware)**

---

## 5. Small-Size Readability Tests (Must Do)

After changes, test the logo at:

* 128×128
* 64×64
* 32×32

At 32×32, it should still be clear that:

* The mascot is a gopher
* The object is some kind of device/controller

If either is unclear, simplify shapes further.

---

## 6. Recommended Variants (Next Step)

Once the main logo is updated, consider creating:

* **Full logo** (gopher + PLC panel + text)
* **Compact logo** (gopher + GoPLC text)
* **Icon mark** (gopher only, square)
* **Monochrome version** (light and dark)

---

## Summary

Priority order:

1. Clarify gopher identity (eye + tooth)
2. Reduce gopher height for balance
3. Optimize for square composition
4. Validate at small sizes

These changes should significantly improve recognition, usability, and polish without altering the core identity of the GoPLC logo.
