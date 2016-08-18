package libs

import (
	"github.com/revel/revel"
	"math"
	"strconv"
)

type PagerResult struct {
	Page     int
	Offset   int
	Limit    int
	Total    int
	MaxPages int
	Current  int
	First    int
	Prev     int
	Next     int
	Last     int
	//	Left     int
	//	Right    int
	Pages []int
	Items []string
	Start int
	End   int
}

type Pager struct {
	Params *revel.Params
	Limit  int
	Left   int
	Right  int
	Items  []string
	Total  int
}

func (p Pager) getPage(max_pages int) int {
	param := p.Params.Query["page"]
	page := 0
	if len(param) > 0 {
		page, _ = strconv.Atoi(param[0])
	}
	if page <= 0 {
		page = 1
	}

	if page > max_pages {
		page = max_pages
	}

	return page
}

func (p Pager) getItems(start int, end int) []string {
	item := p.Items[start:end]

	return item
}

func (p Pager) getStartPosition(offset int) int {
	return offset
}

func (p Pager) getEndPosition(start int) int {
	end := start + p.Limit
	if p.Total < end {
		end = len(p.Items)
	}

	return end
}

func (p Pager) getOffset(page int) int {
	offset := p.Limit * (page - 1)

	return offset
}

func (p Pager) getPrev(page int) int {
	prev := page - 1
	if prev <= 0 {
		prev = 1
	}

	return prev
}

func (p Pager) getNext(page int, max_page int) int {
	next := page + 1
	if next > max_page {
		next = max_page
	}

	return next
}

func (p Pager) getLeftPages(current int, max_pages int) []int {
	var pages []int
	left := p.Left
	for i := current - 1; i > 0 && left >= 1; i-- {
		left--
		p := []int {
			i,
		}
		pages = append(p, pages...)
	}

	return pages
}

func (p Pager) getRightPages(current int, max_pages int, pages []int) []int {
	right := p.Right
	for i := current + 1; i <= max_pages && right >= 1; i++ {
		right--
		pages = append(pages, i)
	}

	return pages
}

func (p Pager) getPages(current int, max_pages int) []int {
	var pages []int
	pages = p.getLeftPages(current, max_pages)
	pages = append(pages, current)
	pages = p.getRightPages(current, max_pages, pages)

	return pages
}

func (p Pager) getMaxPages(total int, limit int) int {
	max_pages := int(math.Ceil(float64(total) / float64(limit)))
	if max_pages <= 0 {
		max_pages = 1
	}

	return max_pages
}

func (p Pager) Result() *PagerResult {
	limit := p.Limit
	total := p.Total
	max_pages := p.getMaxPages(total, limit)
	page := p.getPage(max_pages)
	offset := p.getOffset(page)
	start := p.getStartPosition(offset)
	end := p.getEndPosition(start)
	items := p.getItems(start, end)

	r := &PagerResult{
		Page:     page,
		Limit:    limit,
		Total:    total,
		Offset:   offset,
		Current:  page,
		MaxPages: max_pages,
		First:    1,
		Prev:     p.getPrev(page),
		Next:     p.getNext(page, max_pages),
		Last:     max_pages,
		//		Left:     p.Left,
		//		Right:    p.Right,
		Pages: p.getPages(page, max_pages),
		Items: items,
		Start: start + 1,
		End:   end,
	}

	return r
}
