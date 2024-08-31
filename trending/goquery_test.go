package trending

import (
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// test case for goquery
// Element selector
func TestElementSelector(t *testing.T) {
	html := `
		<body>
			<div>DIV1</div>
			<div>DIV2</div>
			<span>SPAN1</span>
		</body>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		t.Fatalf("create doc err %v", err)
	}

	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("read text: %s", txt)
		if txt != "DIV1" && txt != "DIV2" {
			t.Errorf("invalid text: %s", txt)
		}
	})

	doc.Filter("span").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("text from selection : %s", txt)
		if txt != "SPAN1" {
			t.Errorf("invalid span value")
		}
	})
}

// id selector
func TestIdSelector(t *testing.T) {
	html := `
		<body>
			<div id="div1">DIV1</div>
			<div>DIV2</div>
			<span>SPAN1</span>
		</body>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {

	}

	doc.Find("#div1").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()

		if txt != "DIV1" {
			t.Errorf("ID select incorrect txt: %s", txt)
		}
	})
}

// combine selector:  elecment + ID
func TestElectmentIdSelector(t *testing.T) {
	html := `
		<body>
			<div id="div1">DIV1</div>
			<div>DIV2</div>
			<span>SPAN1</span>
		</body>
	`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	doc.Find("div#div1").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()

		if txt != "DIV1" {
			t.Errorf("read incorrect text, %s", txt)
		}
	})
}

// class selector
func TestClassSelector(t *testing.T) {
	html := `
	<body>
		<div id="div1">DIV1</div>
		<div class="classdiv">DIV2</div>
		<span>SPAN1</span>
	</body>
`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	doc.Find(".classdiv").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()

		if txt != "DIV2" {
			t.Errorf("incorrect text: %s", txt)
		}
	})
}

//element class selector

func TestElementClassSelector(t *testing.T) {
	html := `
		<body>
			<div id="div1">DIV1</div>
			<div class="c1">DIV2</div>
			<span>SPAN1</span>
		</body>
	`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("div.c1").Each(func(i int, s *goquery.Selection) {
		text := s.Text()

		if text != "DIV2" {
			t.Errorf("incorrect text : %s", text)
		}
	})
}

// attribute selector
/*
similar attribute selector
Find("div[lang]")			带有 lang属性的div元素
Find("div[lang=zh]")		带有 lang属性,且value=zh的div元素
Find("div[lang!=zh]")		带有 lang属性,且 value != zh 的div元素
Find("div[lang|zh]")		带有 lang属性,且 value为zh或 以zh-开头的元素
Find("div[lang*=zh]")		带有 lang属性,且 value包含zh这个字符串的 元素
Find("div[lang~=zh]")		带有 lang属性,且 value 包含zh字符串, 单词以 空格分开
Find("div[lang$=zh]")		带有 lang属性,且 value 以zh结尾的div元素, 区分大小写
Find("div[lang^=zh]")		带有 lang属性,且 value 以zh开头的div元素, 区分大小写
*/
func TestAttributeSelector(t *testing.T) {
	html := `
	<body>
		<div id="div1">DIV1</div>
		<div attr="123">DIV2</div>
		<div lang="zh">LANG-DIV</div>
		<div lang="notzh">LANG-notzh</div>
		<div lang="zh-z">LANG-zh-z</div>
		<div lang="czhval">LANG-czhval</div>
		<div lang="czhval t2">LANG-czhval-t2</div>
		<div lang="czh">LANG-czh</div>
		<div lang="zhvalue">LANG-zhvalue</div>
		<span>SPAN1</span>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	doc.Find("div[attr]").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		attrname, ok := s.Attr("attr")

		if ok && "123" != attrname {
			t.Errorf("invalude attribute value: %s", attrname)
		}

		if txt != "DIV2" {
			t.Errorf("select incorrect value: %s", txt)
		}
	})

}

// attribute table test driven
func TestAttributeTable(t *testing.T) {
	html := `
	<body>
		<div id="div1">DIV1</div>
		<div attr="123">DIV2</div>
		<div lang="zh">LANG-DIV</div>
		<div lang="notzh">LANG-notzh</div>
		<div lang="zh-z">LANG-zh-z</div>
		<div lang="czhval">LANG-czhval</div>
		<div lang="czhval t2">LANG-czhval-t2</div>
		<div lang="czh">LANG-czh</div>
		<div lang="zhvalue">LANG-zhvalue</div>
		<span>SPAN1</span>
	</body>
`
	tables := []struct {
		selector   string
		expected   []string
		selectvals []string
	}{
		{
			selector:   "div[lang=zh]",
			expected:   []string{"LANG-DIV"},
			selectvals: []string{},
		},
		{
			selector:   "div[lang!=zh]",
			expected:   []string{"DIV1", "DIV2", "LANG-notzh", "LANG-zh-z", "LANG-czhval", "LANG-czhval-t2", "LANG-czh", "LANG-zhvalue"},
			selectvals: []string{},
		},
		{
			selector:   "div[lang|=zh]",
			expected:   []string{"LANG-DIV", "LANG-zh-z"},
			selectvals: []string{},
		},
		{
			selector:   "div[lang*=zh]",
			expected:   []string{"DIV1", "DIV2", "LANG-DIV", "LANG-notzh", "LANG-zh-z", "LANG-czhval", "LANG-czhval-t2", "LANG-czh", "LANG-zhvalue"},
			selectvals: []string{},
		},
		{
			selector:   "div[lang~=zh]",
			expected:   []string{"DIV1", "DIV2", "LANG-DIV", "LANG-notzh", "LANG-zh-z", "LANG-czhval", "LANG-czhval-t2", "LANG-czh", "LANG-zhvalue"},
			selectvals: []string{},
		},
		{
			selector:   "div[lang$=zh]",
			expected:   []string{"DIV1", "DIV2", "LANG-DIV", "LANG-notzh", "LANG-zh-z", "LANG-czhval", "LANG-czhval-t2", "LANG-czh", "LANG-zhvalue"},
			selectvals: []string{},
		},
		{
			selector:   "div[lang^=zh]",
			expected:   []string{"DIV1", "DIV2", "LANG-DIV", "LANG-notzh", "LANG-zh-z", "LANG-czhval", "LANG-czhval-t2", "LANG-czh", "LANG-zhvalue"},
			selectvals: []string{},
		},
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	attrvals := []string{}
	for _, val := range tables {
		clear(attrvals)
		doc.Find(val.selector).Each(func(i int, s *goquery.Selection) {
			txt := s.Text()
			// attribute
			attrval, ok := s.Attr("lang")
			if ok {
				attrvals = append(attrvals, strings.Trim(attrval, " "))
			}
			val.selectvals = append(val.selectvals, txt)
			if !slices.Contains(val.expected, txt) {
				t.Errorf("%s select incorrect value %s, expected vals: %v", val.selector, txt, val.expected)
			}
		})
		t.Logf("%s select values: %v, attributes: %v", val.selector, val.selectvals, attrvals)
	}
}

// parent>child:  directly children selector
func TestParentChildSelector(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<div>SPAN1</div>
		</span>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	selected := []string{}
	doc.Find("body>div").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		selected = append(selected, txt)
	})
	t.Logf("selected vals: %v", selected)
}

// parent child:  select all children include grandchild

func TestParentChildIncludeGrandChild(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<div>SPAN1</div>
		</span>
	</body>
`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	selected := []string{}
	doc.Find("body div").Each(func(i int, s *goquery.Selection) {

		selected = append(selected, s.Text())
	})
	t.Logf("selected values: %v", selected)
}

// prev+next:  neighbor selector
func TestNeighBorSelect(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<div>SPAN1</div>
		</span>
		<p>P2</p>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	doc.Find("div[lang=zh]+p").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("Select value :%v", txt)
	})
}

// 兄弟选择器: prev~next
func TestBrotherSelector(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<div>SPAN1</div>
		</span>
		<p>P2</p>
	</body>
`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	doc.Find("div[lang=zh]~p").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("brother selector select val: %v", txt)
	})
}

// 内容过滤器:  content selector

func TestContentSelector(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<div>SPAN1</div>
		</span>
		<p>P2</p>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	doc.Find("div:contains(DIV)").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("content select val: %v", txt)
	})
}

// :first-child filter:  筛选出 parent 第一个 child
func TestFirstChild(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<div>SPAN1</div>
		</span>
		<p>P2</p>
	</body>
`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	expect := []string{"DIV1", "SPAN1"}
	doc.Find("div:first-child").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("first child select value: %s", txt)
		if !slices.Contains(expect, txt) {
			t.Errorf("invalid first child : %v", txt)
		}
	})
}

// first-of-type : filter
func TestFirstType(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<p>span first child</p>
			<div>SPAN1</div>
		</span>
		<p>P2</p>
	</body>
`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	expected := []string{"DIV1", "SPAN1"}
	doc.Find("div:first-of-type").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		selectHtml, _ := s.Html()
		t.Logf("select val %s, html: %s", txt, selectHtml)

		if !slices.Contains(expected, txt) {
			t.Errorf("first of type incorrect value : %s", txt)
		}
	})

}

// :last-child and :last-of-type

func TestLastChild(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<p>span first child</p>
			<div>SPAN1</div>
			<div>SPAN2</div>
			<div>SPAN3</div>
		</span>
		<p>P2</p>
	</body>
`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	expected := []string{"SPAN3"}
	doc.Find("div:last-child").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		selectHtml, _ := s.Html()
		t.Logf("select val %s, html: %s", txt, selectHtml)

		if !slices.Contains(expected, txt) {
			t.Errorf("last child incorrect value : %s", txt)
		}
	})

}
func TestLastType(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<p>span first child</p>
			<div>SPAN1</div>
			<div>SPAN2</div>
			<div>SPAN3</div>
			<p>p3</p>
			<div>SPAN4</div>
		</span>
		<p>P2</p>
	</body>
`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	expected := []string{"DIV3", "SPAN4"}
	doc.Find("div:last-of-type").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		selectHtml, _ := s.Html()
		t.Logf("select val %s, html: %s", txt, selectHtml)

		if !slices.Contains(expected, txt) {
			t.Errorf("first of type incorrect value : %s", txt)
		}
	})
}

// :nth-child(n)
func TestNthChild(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<p>span first child</p>
			<div>SPAN1</div>
			<div>SPAN2</div>
			<div>SPAN3</div>
			<p>p3</p>
			<div>SPAN4</div>
		</span>
		<p>P2</p>
	</body>
`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("div:nth-child(3)").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()

		if txt != "DIV2" && txt != "SPAN2" {
			t.Errorf("nth-child invalid value: %s", txt)
		}
	})
}

// :nth-of-type(n)
func TestNthType(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<p>span first child</p>
			<div>SPAN1</div>
			<div>SPAN2</div>
			<div>SPAN3</div>
			<p>p3</p>
			<div>SPAN4</div>
		</span>
		<p>P2</p>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	doc.Find("div:nth-of-type(2)").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()

		if txt != "SPAN2" && txt != "DIV2" {
			t.Errorf("nth-of-type select incorrect value: %s", txt)
		}
	})
}

// :nth-last-child(n)
func TestNthLastChild(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<p>span first child</p>
			<div>SPAN1</div>
			<div>SPAN2</div>
			<div>SPAN3</div>
			<p>p3</p>
			<div>SPAN4</div>
		</span>
		<p>P2</p>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("div:nth-last-child(1)").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()

		if txt != "DIV3" && txt != "SPAN4" {
			t.Errorf("nth-last-child(1) invalide value: %s", txt)
		}
	})
}

// :nth-last-of-type(n)
func TestNthLastOfType(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<p>P1</p>
		<div lang="zh-cn">DIV2</div>
		<div lang="en">DIV3</div>
		<span>
			<p>span first child</p>
			<div>SPAN1</div>
			<div>SPAN2</div>
			<div>SPAN3</div>
			<p>P3</p>
			<div>SPAN4</div>
		</span>
		<p>P2</p>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("p:nth-last-of-type(1)").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()

		if txt != "P2" && txt != "P3" {
			t.Errorf("nth-last-of-type(1) invalide value: %s", txt)
		}
	})
}

// :only-child
func TestOnlyChild(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<span>
			<div>DIV2</div>
		</span>
		<p>P2</p>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("div:only-child").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("only-child : %s", txt)
		if txt != "DIV2" {
			t.Errorf("only-child invalide value: %s", txt)
		}
	})
}

// only-of-type
func TestOnlyOfType(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<span>
			<div>DIV2</div>
		</span>
		<p>P2</p>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("div:only-of-type").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("only-child : %s", txt)
		if txt != "DIV2" && txt != "DIV1" {
			t.Errorf("only-of-type invalide value: %s", txt)
		}
	})
}

// 选择器或运算(,)
func TestSelectOr(t *testing.T) {
	html := `
	<body>
		<div lang="zh">DIV1</div>
		<span>
			<div>DIV2</div>
		</span>
		<div lang="ch">DIV3</div>
		<p>P2</p>
	</body>
`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("div[lang=zh],p").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		t.Logf("or select value: %s", txt)
		if txt != "DIV1" && txt != "P2" {
			t.Errorf("select or incorrect value: %s", txt)
		}
	})
}
