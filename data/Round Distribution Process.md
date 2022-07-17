The Flow Chart (at https://mermaid-js.github.io/mermaid-live-editor/edit#pako:eNqFkt2KwjAQRl8l5FpfoCyC9EdXRYqVvWm9GJtpDZs0JU0qi_XdN-1WUUE2VyHnZPJlkgvNFUPq0VJDfSL7IKuIG_N0i2eyU7ZiBzKdzoif-koIzA2JQYMEg7o5_Ln-IHQL3saa59j5T8vzFriAo0DihGcW3GsGYGCsFgwoTJNc89p4H0c9CwUveV8hUBUYrqrmyY0e3a0y5B9_8ejHVtfONIp8oeYFRza64eAu06TWCKw5IZqRRG_J4i1ZDmSV7hEk2WHL8TyS1UDW6Q6lapFEAsryHmI9wM1jYB9EboXrf9_PI-Tft8ttxn4X7mnGI7rPNFLaShKrxpApiXjFXSr2soOx_v5JH7d7DU8nVKKWwJn7I5d-X0bNCSVm1HNThgVYYTKaVVen2pq5ZCHjRmnqFSAanFCwRiU_VU49oy3epICD-3JytK6_cATO0w

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
    H --> J[Team Review]
    J --> K[Remove Flagged]
    K --> L[Script:<br>Calculate Givbacks]
    L --> |After Review|I[Forum Post - Finished]
    L --> |Add to Sheet|H[Spreadsheet]
