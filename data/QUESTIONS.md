- Can we get some data scientists involved?

- What does the address in the Giveth API mean?

- Many of the projects are totally inactive, would information about usage be interesting?
    - This is the leader board, but there could be more resolution
        - Usage by category
        - Usage over time
        - Usage by community

- Is there an endpoint for getting a list of all projects?

- What is the $NICE token?
    - If a user donates to the Giveth project, they get $NICE in 1:1 of the amount donated
    - Donated money goes into the Giveth Multisig. Leaves the multi-sig to:
        - 1/2 Buys Giv
        - 1/2 Provides Giv/StableCoin pairs
    - Accounting?        

- When I query against the not-eligible-donations and eligible-donations endpoints, I do not get the project's address or id. This makes it very difficult to join data files

- I need a way to determine the dates of the rounds
  - I propose 14 day rounds starting the end of day the first Thursday of January and 14 days thereafter

- I need a clear idea of why some transactions are eligible and some are not.

- Who maintains the purple list - are addresses ever removed from the purple list? What exactly is the purple list?

- How do you select the block number for pricing of Giv? (chifra when End Of Round)

- If I run the same query with identical parameters, I get different results. Simple fix -- add one additional field to the sort of the results: txhash

- The "process" cares about Rounds -- the API interface forces users to enter date ranges
  - the API should allow them to pick a round (easy to produce) and fill in the date ranges 
  - less error prone
  - much easier to reproduce
  - much easier to use
  - all APIs should take 'round'

- Without reproducability, you can't write valid tests

- Need to be able to download "project" information

- Is eligbile plus not-eligible all the transactions that went through the website?

- If a project has two wallets in the address array, which one dominates in the walletAddress field?
    - There are four addresses in the data
        - address
        - admin address
        - addresses::gnosis wallet
        - addresses::mainnet wallet

- What is the difference between walletAddress and adminUser.walletAddress?

- Each of the three queries (eligible, not-eligible, and purpleToVerified) are identical (other than "info"). They should include a field called "type" to make combining the rows easier

- the same txHash in two different files?
jrush@web3:~/D/trueblocks-giveth|main⚡*➤ cat data/summaries/all_donations.csv | cut -d, -f7 | sort | wc
    9620    9620  663720
jrush@web3:~/D/trueblocks-giveth|main⚡*➤ cat data/summaries/all_donations.csv | cut -d, -f7 | sort -u | wc
    9295    9295  641295

- You should change the name from xDAI to gnosis

- Is there an API endpoint to retrieve the givbacks?

- How does the "ok" field in the spreadsheet get re-inserted into the process to create Givbacks?

- Are Givbacks given for donations of non GIV donations (yes)

- How is USD value calculated for the different tokens donated?
