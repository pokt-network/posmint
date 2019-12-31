#LATEST CHANGES (12/12/2019)

##File Changes:

### CREATED / MODIFIED
- account.go
- address.go
- alias.go
- auth/module.go
- auth/util/util.go
- bank/module.go
- CHANGELOG.md
- cliContext.go
- common_test.go
- config.go
- config.yml
- crisis/module.go
- dao_test.go
- dec_coin.go
- ed25519.go
- genesis.go
- go.mod
- go.sum
- hooks.go
- keeper/pool.go
- keeper/validator.go
- keeper/valUtil.go
- keeper/valUtil_test.go
- keybase.go
- keys.go
- lazy_keybase.go
- params.go
- pos/module.go
- query.go
- reward.go
- slash.go
- stdsignmsg.go
- stdtx.go
- supply/module.go
- test_common.go
- tx.go
- tx.go
- txbuilder.go
- types.go
- types/modules/module.go
- types/pool.go
- types/validator.go
- types/valUtil.go
- types/valUtil_test.go
- uint.go
- util_test.go
- valPrevState.go
- valUnstaked.go

### DELETED
- account_retriever.go
- ante.go
- ante_test.go
- gobash.go
- hooks_test.go
- io.go
- knownValues.go
- process.go
- test_cover.sh
- tests/util.go
- types_test.go

## Commit Detail:

1 - Adding Tests to x/pos/types Otto V 12/12/19, 1:32 AM 697674bd

`hooks.go
pool.go
validator.go`

2 - update getValidatorReward to use MustUnmarshalBinaryBare instead of mustUnmarhsalBinaryLengthPrefixed Eduardo Diaz 12/12/19, 11:13 AM 5ae4f995

`reward.go`

3 - Merge branch 'issue-#6' of github.com:pokt-network/posmint into issue-#6 Eduardo Diaz 12/12/19, 11:14 AM c00fb98b

`reward.go
config.go
tmNodeKey.go
keybase.go`

4 - Fixed keybase ed25519 Luis de Leon 12/12/19, 12:13 PM d9dbecaf

`keybase.go`

5 - WIP remove bip39 support Luis de Leon 12/14/19, 11:26 AM b7f696cd

`keybase.go
lazy_keybase.go
types.go`

6 - fixed import export test andrewnguyen22 12/14/19, 11:52 AM 81820790

`keybase.go`

7 - Code cleanup Luis de Leon 12/14/19, 12:08 PM f2d89d4b

`keybase.go
types.go`

8 - Merge pull request #21 from pokt-network/issue-#20 Andrew Nguyen* 12/14/19, 12:10 PM b6f35a1a

`keybase.go
lazy_keybase.go
types.go
reward.go
hooks.go
pool.go
validator.go`

9 - Closes issue #22. (#23) Otto V* 12/18/19, 9:05 AM 0818d20b

`uint.go`

10 - Closes issue #24. (#25) Otto V* 12/18/19, 9:07 AM da0cc0f5

`dec_coin.go`

11 - set default default proposer reward percentage; use common_test utility Eduardo Diaz* 12/17/19, 12:22 PM 3cd7f55e

`reward.go
params.go`

12 - format project Eduardo Diaz* 12/17/19, 12:22 PM 714ebce8

`account.go
alias.go
config.go
ed25519.go`

13 - Got rid of gas Luis de Leon 12/12/19, 3:40 PM 4a3b99e8

`stdsignmsg.go
stdtx.go
test_common.go
txbuilder.go`

14 - Removed pending Stdfee and Gas related calls Luis de Leon 12/18/19, 3:51 PM 37c072e0

`util.go
alias.go`

15 - Merge pull request #26 from pokt-network/issue-#6 Luis de Leon* 12/18/19, 4:11 PM c7b9ef28

`config.go
ed25519.go
account.go
stdsignmsg.go
stdtx.go
test_common.go
txbuilder.go
util.go
alias.go
ante.go
params.go`

16 - Added circleci configuration Luis de Leon 12/18/19, 4:33 PM 5b27e086

`config.yml`

17 - updated appmodules for setting keybase and tendermint node andrewnguyen22* 12/20/19, 10:16 AM 84ed6e43

`module.go
module.go
module.go
module.go
module.go
module.go`

18 - patch keeper mint to use moduleStakedPool instead of modulename Eduardo Diaz* 12/20/19, 7:05 PM cb4daa0f

`reward.go`

19 - fixed minting, registered necessary codecs andrewnguyen22* 12/21/19, 2:57 PM ce5c9e11

`reward.go`

20 - fixed genesis bug andrewnguyen22* 12/21/19, 4:27 PM 8cb8a68e

`module.go
module.go
module.go
module.go
module.go
module.go`

21 - fixes signing info bug andrewnguyen22* 12/21/19, 5:05 PM 403df57d

genesis.go

22 - removing bech32 Otto V* 12/20/19, 11:50 AM c72cefef

`address.go
account.go
stdtx.go
valUtil.go`

23- changing main prefix to "pocket" Otto V* 12/23/19, 11:25 AM 4a0061d2

`address.go`

24 - bech32 Public Key of validator to Hex andrewnguyen22* 12/23/19, 2:25 PM 3ee77c13

`valUtil.go
valUtil.go
genesis.go
address.go`

25 - Merge branch 'staging' of https://github.com/pokt-network/posmint into issue-#28 Otto V 12/23/19, 4:54 PM 07a3deda

`module.go
address.go
account.go
stdtx.go
stdtx_test.go
module.go
module.go
module.go
reward.go
valUtil.go
valUtil.go
genesis.go
module.go
module.go`

26 - update POSMint without keybase and tendermint node andrewnguyen22* 12/26/19, 8:46 PM f3ce8a09

`module.go
cliContext.go
module.go
module.go
module.go
tx.go
module.go
query.go
tx.go
module.go`

27 - Merge branch 'staging' of https://github.com/pokt-network/posmint into issue-#28 Otto V 12/27/19, 5:36 PM 26d85227

`module.go
cliContext.go
module.go
module.go
module.go
tx.go
module.go
query.go
tx.go
module.go
valUtil.go`

28 - Adding CHANGELOG.md Otto V 12/30/19, 10:19 AM e0c7caff

`CHANGELOG.md`

29 - Adding CHANGELOG.md Otto V 12/30/19, 10:22 AM 43422470

`CHANGELOG.md`

30 - Updating some tests, running go fmt Otto V 12/27/19, 5:36 PM f0652961

`valUtil.go`

31- Merge branch 'issue-#28' into staging Otto V 12/30/19, 10:42 AM b29f24b4

`valUtil.go`

32 - add dao tests Eduardo Diaz 12/23/19, 12:26 PM 6c0c48c9

`reward.go`

33 - test keeper pool staking, unstaking and burn Eduardo Diaz 12/24/19, 5:06 PM 557ef871

`pool.go`

34 - remove unused private getValsFromPrevState method from pos module keeper Eduardo Diaz 12/27/19, 3:30 PM 8a2355e6

`valPrevState.go`

35 - validaotr stake tests Eduardo Diaz 12/28/19, 1:36 AM a6fe8c64

`valPrevState.go`

36 - test validator unstaking methods Eduardo Diaz 12/28/19, 6:14 PM 124abf03

`valUnstaked.go`

37 - test validator utilities Eduardo Diaz 12/29/19, 3:54 AM efd7458c

`valUtil.go`

38 - test keeper slash methods Eduardo Diaz 12/30/19, 2:49 AM 41d096dd

`reward.go
slash.go`

39 - fixed burn and validator bug andrewnguyen22* 12/30/19, 10:42 AM 969a1331

`reward.go
slash.go
validator.go
valPrevState.go
keys.go`

40 - fixed base proposer percentage bug andrewnguyen22* 12/30/19, 12:01 PM f4b42c8a

`reward.go`

41 - fixed expected from tests, gofmt andrewnguyen22* 12/30/19, 4:24 PM 8daefbb8

`module.go
module.go
module.go`



