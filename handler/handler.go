package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

type Borrowing struct {
	User string `json:"user"`
	Item string `json:"item"`
}

type Returning struct {
	Item string `json:"item"`
}

func (b1 Borrowing) equal(b2 Borrowing) bool {
	return b1.User == b2.User && b1.Item == b2.Item
}

func (b1 Borrowing) equalItem(b2 Borrowing) bool {
	return b1.Item == b2.Item
}

func Borrow() echo.HandlerFunc {
	return func(c echo.Context) error {
		var b Borrowing
		bytes, err := ioutil.ReadAll(c.Request().Body)

		if err != nil {
			c.Error(err)
		}

		err = json.Unmarshal(bytes, &b)
		if err != nil {
			c.Error(err)
		}

		current, err := readItems()
		if err != nil {
			c.Error(err)
		}

		result := merge(current, b)

		err = writeItems(result)
		if err != nil {
			c.Error(err)
		}

		return c.String(200, "ok")
	}
}

func Return() echo.HandlerFunc {
	return func(c echo.Context) error {
		bytes, err := ioutil.ReadAll(c.Request().Body)

		if err != nil {
			c.Error(err)
		}

		var r Returning

		err = json.Unmarshal(bytes, &r)
		if err != nil {
			c.Error(err)
		}

		current, err := readItems()

		result := delete(current, r)

		err = writeItems(result)
		if err != nil {
			c.Error(err)
		}

		if err != nil {
			c.Error(err)
		}

		return c.String(200, "ok")
	}
}

func delete(bs []Borrowing, r Returning) []Borrowing {
	var result = make([]Borrowing, 0)

	for _, target := range bs {
		if !(target.Item == r.Item) {
			result = append(result, target)
		}
	}

	return result
}

func merge(bs []Borrowing, b Borrowing) []Borrowing {
	var result = make([]Borrowing, 0)

	for _, target := range bs {
		if !target.equalItem(b) {
			result = append(result, target)
		}
	}

	result = append(result, b)
	return result
}

func Items() echo.HandlerFunc {
	return func(c echo.Context) error {
		items, err := readItems()

		if err != nil {
			c.Error(err)
		}

		return c.JSON(http.StatusOK, items)
	}
}

const pathname = "./tmp/data.json"

func writeItems(b []Borrowing) error {
	bytes, err := json.Marshal(b)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pathname, bytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

func readItems() ([]Borrowing, error) {
	file, err := os.OpenFile(pathname, os.O_RDONLY, 0666)

	if os.IsNotExist(err) {
		return []Borrowing{}, nil
	} else if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	var items []Borrowing
	if err := json.Unmarshal(bytes, &items); err != nil {
		return nil, err
	}

	return items, nil
}
