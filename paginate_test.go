package paginate

import (
	"fmt"
	"testing"
)

func TestGetPagingOptions(t *testing.T){
	in := &PageOptions{
		PageSize:-1,
		PageNumber:3,
		//CurrPageNumber:1,
	}
	offset, limit,symbol := GetPagingOptions(in)
	fmt.Println(offset,limit,symbol)
}