# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] 2020-05-31

### Changed
- Remove need for keybase BuildAndSign. [#200]

## [Unreleased] 2020-01-15

### Added
- Created alternate func HexAddressPubKeyAmino that returns the amino encoded bytes [#69]

### Changed
- Standardize the Address type to be used instead of AccAddress , ValAddress and ConsAddress.[#65]
- Refactor Validators and Related. [#65]
- Changed HexAddressPubKey to use actual bytes and not amino encoded bytes.[#69]
- Refactor Test (Redundant Values in func after Address type) [#67]
- Refactor Hooks (Redundant Values in func after Address type) [#67]
- Refactored stdtx to have one message and one signature [#186]

### Removed
- Clean Remaining Bech32 [#65]

## [Unreleased] 2019-12-30

### Added
- Add Unit Tests to x/pos/types [#28]
- Add setTendermintNode() and setKeybase() to app modules [#31]
- Setup CI [#27]
- Add unit tests to the POS package [#16]
- Add unit tests for Keybase and Lazy Keybase [#14]

### Changed
- Updated latest changes to this document [#47]
- Changing Bech32 to Hex [#35]
- Refactor client [#11]
- Refactor crypto/keys [#10]
- Replace Distribution/Staking/Slashing/Minting with single POS module [#1]

### Removed
- Remove keybase and tendermint node from app module [#37]
- Remove BIP 39 support for key creation/import [#20]
- Removed Gas Usage [#18]
- Clean cosmos-sdk from project [#6]
- Remove SIMAPP [#4]
- Remove CLI and RPC functionality [#3]

### Fixed
- Fix Multiple Tests related to x/pos/keepers [#43,#44,#48,#49,#52]
- Fix ValidateDoubleSign tests [#41]
- Init Genesis never called if genesis.json isn't completely filled out [#33]
- Undesirable "negative coin amount" in DecCoin's Sub[#24]
- Uint.LTE() BUG [#22]


[Unreleased]: https://github.com/pokt-network/posmint/compare/master...staging