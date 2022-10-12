package data

func GetCategories() map[string][]string {
	projects := GetProjects()
	cats := make(map[string][]string)
	for _, project := range projects {
		for _, c := range project.Categories {
			cats[c.Name] = append(cats[c.Name], project.Title)
		}
	}
	return cats
}
