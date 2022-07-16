The Flow Chart (at https://mermaid-js.github.io/mermaid-live-editor/edit#pako:eNqFkl1rwjAUhv9KyPX8A2UI0g-dihQru2m9ODanNSxpSppUhvW_L-2q2IEsVyHv856cryvNFUPq0VJDfSaHIKuIO4t0hxeyV7ZiRzKbzYmf-koIzA2JQYMEg7o5_rL-AHRL3saa59j5k-dFC1zASSBxwFQLHjEDMDBGCwYpTJNc89p47yc9DwUveR8hUBUYrqpmwkbP7E4Z8g-_fOZjq2tHGkU-UfOCIxvZcGBXaVJrBNacEc2oRC-V5UtlNTajcH0je2w5XrqPNFLaShKrxpAZiXjFnYVNHOv0gCBHx6isB2WT7lGqFkkkoCwfts0gbp9L9EHkVriJ9RM4Qf51b8d2TIqxvv6kT7f7mzx9oxK1BM7cjlx7X0bNGSVm1HNXhgVYYTKaVTeH2pq5f0LGjdLUK0A0-EbBGpV8Vzn1jLZ4hwIObuXkSN1-AMiPzs8)

graph TD
    A[New Round] --> C[Collect Paramaters]
    C --> |GivPrice|C
    C --> |Available Giv|C
    C --> D[Collect Data]
    D --> E[Script:<br>Eligible Donations]
    D --> F[Script:<br>Not Eligible Donations]
    D --> G[Script:<br>Purple to Verified]
    E --> H[Spreadsheet]
    F --> H[Spreadsheet]
    G --> H[Spreadsheet]
    H --> |After Review|I[Forum Post - Finished]
    H --> J[Team Review]
    J --> K[Remove Flagged]
    K --> L[Script:<br>Calculate Givbacks]
    L --> |Add to Sheet|H[Spreadsheet]