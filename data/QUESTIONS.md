# Outstanding / Open Questions

### Quotes and Links from Giveth.io Website

- Circulating donations raised by other means. Only “first touch” donations count for GIVbacks. If a project receives funding from a donor and is found to be circulating these donations within the Giveth platform to receive GIVbacks, they will be disqualified. For example, a project should not be sending fiat donations received elsewhere back to their donors and asking them to donate on Giveth with crypto.
-       https://docs.giveth.io/giveconomy/givbacks/#disqualifying-factors-for-the-givbacks-program

- List of eligible Tokens: https://forum.giveth.io/t/givbacks-token-list/253

### Questions In General

- You should change the name from xDAI to gnosis

- If I run the same query with identical parameters, I get different results. Simple fix -- add one additional field to the sort of the results: txhash

- Without reproducability, you can't write valid tests

### Questions About Website Data

- When I query against the not-eligible-donations and eligible-donations endpoints, I do not get the project's address or id. This makes it very difficult to join data files.

- What is the distinction between eligible transactions and not-eligible transasctions?

- How is the purple list maintained 
  - are addresses ever removed from the purple list?
  - Why are addresses added to the purple list?

- How os block number selected for the GivBack calc? (chifra when End Of Round)

- I need a way to determine the dates of the rounds
  - I propose 14 day rounds starting the end of day the first Thursday of January and 14 days thereafter

- The users care about Rounds
  - the API should allow them to pick a round (easy to produce) and fill in the date ranges 
  - less error prone
  - much easier to reproduce
  - much easier to use
  - all APIs should take 'round'

- Need to be able to download the project information. How does one download the list of all projects along with all the project's data?

- Is "eligbile" plus "not-eligible" all the transactions that went through the website?

- If a project has two wallets in the address array, which one dominates in the walletAddress field?

- What is the difference between walletAddress, adminUser.walletAddress, and the numerous addresses in the addresses field?

- Each of the three queries (eligible, not-eligible, and purpleToVerified) should proudce identical data, but they do not (one of them has an "info" field). 

- The data produced by these three endpoints should have a "type" field to make combining them easier. Types: "eligible", "not-eligible", "purpleList-to-verified-projects"

- Is it ever the case that there is two different files with the same transaction hash?

- If I combine all `not-eligible` and all `eligible` transactions, are these two things true:
  - That is a list of all donations to the website since inception
  - There are no duplicates (I think the answer to this is no)

### Possible Article Topics

#### Tracing Suspicious Donations on Giveth

1. Thanks Giveth
2. Link to the original proposal
3. Talk about the Giveth platform (one paragraph and a pointer)
4. Talk about TrueBlocks (decentralized - speed -- accuracy -- no timeouts)
5. Data Definitions
   1. Projects
      1. address,grantId,name,tag,active,core
   2. Donations
      1. type,round,amount,currency,createdAt,valueUsd,giverAddress,txHash,network,source,giverName,giverEmail,projectLink
   3. Givbacks
      1. type,round,givDistributed,givFactor,givPrice,givbackUsdValue,giverAddress,giverEmail,giverName,totalDonationsUsdValue,givback,share
   4. Rounds
      1. id,startDate,endDate,gvailable,price
6. Define the history of the token
7. Define the problem
8. Describe the solution

#### Studying the Giveth Token with TrueBlocks

chifra names Giveth ERC20
Address is:

chifra list <address> --verbose 221,100 appearances
Chart of activity by data by day

Custom display string to include date
chifra export <address> --articulate | cut out input data
Chart of behavior by day

How long to things take? 7:42 - 8:49