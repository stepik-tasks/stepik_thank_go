package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// начало решения

type Department struct {
	XMLName xml.Name `xml:"department"`

	Code      string     `xml:"code"`
	Employees []Employee `xml:"employees>employee"`
}

type Employee struct {
	XMLName xml.Name `xml:"employee"`

	Id     string `xml:"id,attr"`
	Name   string `xml:"name"`
	City   string `xml:"city"`
	Salary string `xml:"salary"`
}

func (emp *Employee) Slice(dep *Department) []string {
	salary := "0"
	id := "0"

	if emp.Salary != "" {
		salary = emp.Salary
	}
	if emp.Id != "" {
		id = emp.Id
	}

	return []string{id, emp.Name, emp.City, dep.Code, salary}
}

func CsvHeader() []string {
	return []string{"id", "name", "city", "department", "salary"}
}

// ConvertEmployees преобразует XML-документ с информацией об организации
// в плоский CSV-документ с информацией о сотрудниках
func ConvertEmployees(outCSV io.Writer, inXML io.Reader) error {
	// Чтение department по частям

	decoder := xml.NewDecoder(inXML)
	writer := csv.NewWriter(outCSV)

	var err error

	err = writer.Write(CsvHeader())
	if err != nil {
		return err
	}

	var valid bool
	valid = false

	for {
		var d Department
		var tok xml.Token

		// берем следующий токен
		tok, err = decoder.Token()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch tp := tok.(type) {
		case xml.StartElement:
			if tp.Name.Local == "organization" {
				valid = true
			}
			if tp.Name.Local == "department" {
				// Декодирование department
				err = decoder.DecodeElement(&d, &tp)

				if err != nil {
					return err
				}

				for _, e := range d.Employees {
					err = writer.Write(e.Slice(&d))
					if err != nil {
						return err
					}
				}
			}
		}
	}

	if !valid {
		return fmt.Errorf("wrong xml format")
	}

	writer.Flush()
	if writerErr := writer.Error(); writerErr != nil {
		return writerErr
	}

	if err == io.EOF {
		return nil
	}

	return err
}

// конец решения

func main() {
	src := `<organization>
    <department>
        <code>hr</code>
        <employees>
            <employee id="11">
                <name>Дарья</name>
            </employee>
            <employee id="12">
                <name>Борис</name>
                <city>Самара</city>
                <salary>78</salary>
            </employee>
        </employees>
    </department>
    <department>
        <code>it</code>
        <employees>
            <employee id="21">
                <name>Елена</name>
                <city>Самара</city>
                <salary>84</salary>
            </employee>
        </employees>
    </department>
</organization>`

	in := strings.NewReader(src)

	out := os.Stdout
	err := ConvertEmployees(out, in)

	fmt.Println("error: ", err)
	/*
		id,name,city,department,salary
		11,Дарья,Самара,hr,70
		12,Борис,Самара,hr,78
		21,Елена,Самара,it,84
	*/

	/*
		id,name,city,department,salary
		11,Дарья,Самара,hr,70
		12,Борис,Самара,hr,78
		21,Елена,Самара,it,84
	*/
}
