# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] 2019-12-30

### Added
- Add Unit Tests to x/pos/types [#28]
- Add setTendermintNode() and setKeybase() to app modules [#31]
- Setup CI [#27]
- Add unit tests to the POS package [#16]
- Add unit tests for Keybase and Lazy Keybase [#14]

### Changed
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
- Init Genesis never called if genesis.json isn't completely filled out [#33]
- Undesirable "negative coin amount" in DecCoin's Sub[#24]
- Uint.LTE() BUG [#22]


[Unreleased]: https://github.com/pokt-network/posmint/compare/master...staging