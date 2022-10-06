package data

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/ethereum/go-ethereum/common"
)

func GetProjects() (projects []Project) {
	filepath.Walk(DataFolder()+"raw/", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			fp, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fp.Close()
			b, _ := io.ReadAll(fp)
			w := Wrapper{}
			json.Unmarshal(b, &w)
			p := w.Data.Project
			if len(p.Id) > 0 && validate.IsValidAddress(p.WalletAddress) {
				p.WalletAddress = strings.ToLower(p.WalletAddress)
				p.Title = strings.Trim(strings.Replace(p.Title, "\n", " ", -1), " ")
				for i := 0; i < len(p.Addresses); i++ {
					p.Addresses[i].Address = strings.ToLower(p.Addresses[i].Address)
				}
				projects = append(projects, p)
			}
		}
		return nil
	})
	return projects
}

type SimpleProject struct {
	Address    string `json:"address"`
	GrantId    string `json:"grantId"`
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	Active     bool   `json:"active"`
	Core       bool   `json:"core"`
	Categories string `json:"categories"`
}

var namesMap *names.NamesMap

func GetChifraName(addrIn string) (*names.Name, error) {
	if namesMap == nil {
		nM, err := names.LoadNamesMap("mainnet")
		if err != nil {
			return nil, err
		}
		namesMap = &nM
	}
	a := common.HexToAddress(addrIn)
	name := (*namesMap)[a]
	return &name, nil
}

func ToSimpleProject(p *Project) SimpleProject {
	isCore := strings.HasPrefix(p.Id, "core-")
	cats := ""
	for _, c := range p.Categories {
		if len(cats) > 0 {
			cats += " "
		}
		cats += c.Name
	}

	return SimpleProject{
		Address:    p.WalletAddress,
		GrantId:    p.Id,
		Name:       p.Title,
		Tag:        "31-Projects:Giveth",
		Active:     p.Listed,
		Core:       isCore,
		Categories: cats,
	}
}

type Project struct {
	Id                  string           `json:"id"`
	Title               string           `json:"title"`
	Balance             float64          `json:"balance"`
	Image               string           `json:"image"`
	Slug                string           `json:"slug"`
	SlugHistory         []string         `json:"slugHistory"`
	CreationDate        string           `json:"creationDate"`
	UpdatedAt           string           `json:"updatedAt"`
	Admin               string           `json:"admin"`
	Description         string           `json:"description"`
	WalletAddress       string           `json:"walletAddress"`
	ImpactLocation      string           `json:"impactLocation"`
	QualityScore        int              `json:"qualityScore"`
	Verified            bool             `json:"verified"`
	TraceCampaignId     *string          `json:"traceCampaignId"`
	Listed              bool             `json:"listed"`
	GivingBlocksId      *string          `json:"givingBlocksId"`
	Status              Status           `json:"status"`
	Categories          []Category       `json:"categories"`
	Reaction            *Reaction        `json:"reaction"`
	AdminUser           User             `json:"adminUser"`
	Organization        Organization     `json:"organization"`
	Addresses           []NetworkAddress `json:"addresses"`
	TotalReactions      int              `json:"totalReactions"`
	TotalDonations      int              `json:"totalDonations"`
	TotalTraceDonations int              `json:"totalTraceDonations"`
}

func (p Project) String() string {
	b, _ := json.MarshalIndent(p, "", "  ")
	return string(b)
}

type NetworkAddress struct {
	Address     string `json:"address"`
	IsRecipient bool   `json:"isRecipient"`
	NetworkId   int    `json:"networkId"`
}

type Organization struct {
	Name                string `json:"name"`
	Label               string `json:"label"`
	SupportCustomTokens bool   `json:"supportCustomTokens"`
}

type User struct {
	Id            int     `json:"id"`
	Email         *string `json:"email"`
	FirstName     string  `json:"firstName"`
	WalletAddress string  `json:"walletAddress"`
}

type Reaction struct {
	Id int `json:"id"`
}

type Category struct {
	Name string `json:"name"`
}

type Status struct {
	Id          int    `json:"id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DataWrapper struct {
	Project Project `json:"projectById"`
}

type Wrapper struct {
	Data DataWrapper `json:"data"`
}
