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

type StringCounter struct {
	Key   string
	Count int
}

type CategoryCounter struct {
	StringCounter
}

func NewCategoryCounter(key string, count int) CategoryCounter {
	return CategoryCounter{StringCounter: StringCounter{Key: key, Count: count}}
}
