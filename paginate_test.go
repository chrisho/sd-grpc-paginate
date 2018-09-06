package paginate

import (
	"fmt"
	"testing"
)

func TestGetPagingOptions(t *testing.T){
	in := &PageOptions{
		PageSize:-1,
		PageNumber:4,
		//CurrPageNumber:1,
		//SortValue:2,
	}
	offset, limit,symbol := GetPagingOptions(in)
	fmt.Println(offset,limit,symbol)
}