package textmail

import (
  "bytes"
  "fmt"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/xml"
  "strings"
)

// Set of whitespace characters.
var whitespace = map[byte]bool{
  ' ':  true,
  '\n': true,
  '\r': true,
  '\t': true,
  '\v': true,
}

// Set of elements we ignore.
var ignore = map[string]bool{
  "audio":  true,
  "head":   true,
  "script": true,
  "track":  true,
  "video":  true,
}

// Set of HTML block elements.
var block = map[string]bool{
  "address":    true,
  "article":    true,
  "aside":      true,
  "blockquote": true,
  "body":       true,
  "canvas":     true,
  "del":        true,
  "div":        true,
  "dl":         true,
  "fieldset":   true,
  "figcaption": true,
  "figure":     true,
  "footer":     true,
  "form":       true,
  "header":     true,
  "hgroup":     true,
  "hr":         true,
  "ins":        true,
  "menu":       true,
  "noscript":   true,
  "ol":         true,
  "output":     true,
  "p":          true,
  "section":    true,
  "table":      true,
  "td":         true,
  "tfoot":      true,
  "th":         true,
  "thead":      true,
  "tr":         true,
  "ul":         true,
}

// HTML heading elements.
var heading = map[string]bool{
  "h1": true,
  "h2": true,
  "h3": true,
  "h4": true,
  "h5": true,
  "h6": true,
}

// HTML elements we will wrap with /
var italic = map[string]bool{
  "em": true,
  "i":  true,
}

// HTML elements we will wrap with *
var bold = map[string]bool{
  "b":      true,
  "strong": true,
}

type Formatter struct {
  buf   bytes.Buffer // Buffer
  links []string     // Collected urls
}

// Converts HTML into nice formatted plain text.
func Format(src []byte) ([]byte, error) {
  doc, err := gokogiri.ParseHtml(src)

  if err != nil {
    return []byte{}, err
  }
  defer doc.Free()

  f := new(Formatter)
  f.walk(doc.Node)

  for i, link := range f.links {
    f.buf.WriteString(fmt.Sprintf("[%d] %s\n", i, link))
  }

  return f.buf.Bytes(), nil
}

// Walks through the documents elements and populates the buffer.
func (self *Formatter) walk(node xml.Node) {
  for c := node.FirstChild(); c != nil; c = c.NextSibling() {
    self.walk(c)
  }

  if node.NodeType() == xml.XML_ELEMENT_NODE {
    self.handleNode(node)
  }
}

// Formats the content of inline elements and writes the content of block
// elements to the buffer.
func (self *Formatter) handleNode(node xml.Node) {
  name := node.Name()

  switch {
  case ignore[name]:
    // Remove ignored elements.
    node.SetContent("")
  case name == "pre":
    // Treat pre elements as code blocks.
    self.writeCodeBlock(node)
  case heading[name]:
    // Headings are prefixed with "# ".
    self.writeBlock(node, "# ")
  case name == "li":
    // List items are prefixed with "- ".
    self.writeBlock(node, "- ")
  case name == "br":
    // Preserve explicit line breaks.
    node.SetContent("\n")
  case italic[name]:
    // Wrap italic elements with /.
    node.SetContent("/" + node.Content() + "/")
  case bold[name]:
    // Wrap bold elements with *.
    node.SetContent("*" + node.Content() + "*")
  case name == "img":
    // Collect the src of images and replace them with (alt)[url index]
    alt, src := node.Attr("alt"), node.Attr("src")

    if len(alt) > 0 && len(src) > 0 {
      node.SetContent(fmt.Sprintf("(%s)[%d]", alt, len(self.links)))
      self.links = append(self.links, src)
    }
  case name == "a":
    // Collect the href and and the url index.
    href, content := node.Attr("href"), node.Content()

    if len(href) > 0 && len(content) > 0 {
      node.SetContent(fmt.Sprintf("%s[%d]", content, len(self.links)))
      self.links = append(self.links, href)
    }
  case block[name]:
    // Write the content of block elements to the buffer.
    self.writeBlock(node, "")
  }
}

// Writes text blocks to the buffer.
func (self *Formatter) writeBlock(node xml.Node, prefix string) {
  block := []byte(strings.TrimSpace(node.Content()))
  node.SetContent("")

  if len(block) == 0 {
    return
  }
  // Position of last space, line break and max length.
  sp, br, max := 0, 0, 79-len(prefix)
  self.buf.WriteString(prefix)

  for i, c := range block {
    // Break line if exceeded max length and the position of the last space
    // is greater than the position of the last line break. Don't break very
    // long words.
    if i-br > max && sp > br {
      self.buf.WriteByte('\n')
      br = sp
      // Only the first line is prefixed.
      for j := 0; j < len(prefix); j++ {
        self.buf.WriteByte(' ')
      }
    }
    if whitespace[c] {
      // The last character was a space, so ignore this one.
      if sp == i {
        sp++
        br++
        continue
      }
      // Write the last word to the buffer, append a space and update
      // the position of the last space.
      if sp > br {
        self.buf.WriteByte(' ')
      }
      self.buf.Write(block[sp:i])
      sp = i + 1
    }
  }

  // Write the last word to the buffer.
  if sp < len(block) {
    if sp > br {
      self.buf.WriteByte(' ')
    }
    self.buf.Write(block[sp:])
  }

  // Close block with 2 breaks.
  self.buf.Write([]byte{'\n', '\n'})
}

// Writes code blocks to the buffer.
func (self *Formatter) writeCodeBlock(node xml.Node) {
  block := []byte(strings.Trim(node.Content(), "\n\r\v"))
  node.SetContent("")

  if len(block) == 0 {
    return
  }
  self.buf.Write(block)
  self.buf.Write([]byte{'\n', '\n'})
}
