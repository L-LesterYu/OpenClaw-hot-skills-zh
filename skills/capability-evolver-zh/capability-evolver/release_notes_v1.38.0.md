## v1.38.0 - Harden Gene Validation Rules

### Breaking Changes (Gene Library)

- `genes.json` bumped to version 2. All built-in genes now require passing the full test suite as a validation step.
- Empty `validation` arrays now score **0.5** instead of **1.0** in PRM, penalizing genes that skip validation.

### Validation Hardening

- **`validate-modules.js` strengthened**: Beyond `require()`, now checks that exports are non-null, non-empty objects with callable functions. Previously, a module exporting `{}` or `null` would pass.
- **New `validate-suite.js`**: Runs the project's test suite (`node --test`) as a Gene validation command. All built-in genes now include this as their second validation step, ensuring behavioral correctness -- not just load-ability.
- **Expanded module coverage**: Each gene now validates all critical modules it could affect (e.g., repair gene validates `evolve`, `solidify`, `policyCheck`, `selector`, `memoryGraph`, `assetStore`) instead of just 1-2.

### PRM Scoring Fix

- `computeProcessScores` Phase 6: genes with empty `validation` arrays previously received a perfect 1.0 validation score. Now they receive 0.5, incentivizing every gene to define at least one validation command.

### Distiller/PolicyCheck Consistency

- `skillDistiller.validateSynthesizedGene()` now uses `policyCheck.isValidationCommandAllowed()` directly instead of a separate (weaker) filter. This prevents the scenario where a distilled gene saves a `node -e "..."` validation command that would be BLOCKED at solidify time.

### Tests

- 266 tests pass, 0 failures
- 24 new test cases covering `computeProcessScores` empty/partial/full validation scoring, `isValidationCommandAllowed` for `--eval`/`-p`/`--print`/`$()`, and distiller `node -e` filtering consistency

**Full Changelog**: v1.37.0...v1.38.0
