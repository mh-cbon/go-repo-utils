package commit

import (
  "time"
  "sort"
)

type Commit struct {
  Revision  string `json:"revision"`
  Author    string `json:"author,omitempty"`
  Email     string `json:"email,omitempty"`
  Date      string `json:"date,omitempty"`
  Message   string `json:"message,omitempty"`
}

func (c Commit) GetDate() *time.Time {
  if c.Date=="" {
    return nil
  }
  d, err := time.Parse("Mon 2006-01-02 15:04:05 -0700", c.Date)
  if err==nil {
    return &d
  }
  d, err = time.Parse("Mon Jan 02 15:04:05 2006 -0700", c.Date)
  if err==nil {
    return &d
  }
  d, err = time.Parse("2006-01-02 15:04:05 -0700", c.Date)
  if err==nil {
    return &d
  }
  return nil
}

type Commits []Commit
type CommitsAsc Commits
type CommitsDesc Commits

func (c Commits) OrderByDate(dir string) {
  if dir=="" {
    dir = "ASC"
  }
  if dir=="ASC" {
    sort.Sort(CommitsAsc(c))
  } else if dir=="DESC" {
    sort.Sort(CommitsDesc(c))
  }
}
func (c Commits) Reverse() {
	for i, j := 0, len(c)-1; i < j; i, j = i+1, j-1 {
		c[i], c[j] = c[j], c[i]
	}
}

// sort utils.
func (s CommitsAsc) Len() int {
  return len(s)
}
func (s CommitsAsc) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}
func (s CommitsAsc) Less(i, j int) bool {
  d1 := s[i].GetDate()
  if d1==nil {
    return false
  }
  d2 := s[j].GetDate()
  if d2==nil {
    return false
  }
  return d2.After(*d1)
}
func (s CommitsDesc) Len() int {
  return len(s)
}
func (s CommitsDesc) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}
func (s CommitsDesc) Less(i, j int) bool {
  d1 := s[i].GetDate()
  if d1==nil {
    return false
  }
  d2 := s[j].GetDate()
  if d2==nil {
    return false
  }
  return d1.After(*d2)
}
