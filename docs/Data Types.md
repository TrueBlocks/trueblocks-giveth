The donation data type

```[hash]
type Donation struct {
  string  type
  string  round
  float64 amount
  string  currency
  string  createdAt
  float64 valueUsd
  string  giverAddress
  string  txHash
  string  source
  string  giverName,omitempty
  string  giverEmail,omitempty
  string  projectlink
}
```

```[bash]
type Givback struct {
  string  type
  string  round
  float64 givDistributed
  float64 givFactor
  float64 givPrice
  float64 givbackUsdValue
  string  giverAddress
  string  giverEmail
  string  giverName
  float64 totalDonationsUsdValue
  float64 givback
  float64 share
}
```

```[bash]
type SimpleProject struct {
  string address
  string grantId
  string name
  string tag
  bool   active
  bool   core
  string categories
}
```

```[bash]
type Project struct {
  string           id
  string           title
  float64          balance
  string           image
  string           slug
  []string         slugHistory
  string           creationDate
  string           updatedAt
  string           admin
  string           description
  string           walletAddress
  string           impactLocation
  int              qualityScore
  bool             verified
  *string          traceCampaignId
  bool             listed
  *string          givingBlocksId
  Status {
    int    id
    string symbol
    string name
    string description
  } status
  []Category {
    string name
  } categories
  *Reaction {
    int id
  } reaction
  User {
    int     id
    *string email
    string  firstName
    string  walletAddress
  } adminUser
  Organization {
    string name
    string label
    bool   supportCustomTokens
  } organization
  []NetworkAddress {
    string address
    bool   isRecipient
    int    networkId
  } addresses
  int totalReactions
  int totalDonations
  int totalTraceDonations
}
```

```[bash]
type Round struct {
  int      id
  DateTime startDate
  DateTime endDate
  int      available
  float64  price
}
```
