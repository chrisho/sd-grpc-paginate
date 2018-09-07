package paginate;

import (
	"strings"
	"github.com/chrisho/mosquito/utils"
	"reflect"
	"github.com/jinzhu/gorm"
)

const (
	PagingSize      = 10
	PagingByPrimary = iota
	PagingByNumber
)

const SortFieldSuffix = "_sort"

func GetPaginate(db *gorm.DB, tname string,in *PageOptions)*gorm.DB{
	offset, limit,symbol := GetPagingOptions(in)
	field := tname +"." + in.SortField
	if in.SortValue > 0 {
		db = db.Where(field + symbol + "?",in.SortValue)
	}
	return db.Order(field + " " + in.SortFieldTo).Offset(offset).Limit(limit)
}
// 获取分页选项
func GetPagingOptions(in *PageOptions) (offset, limit int32, symbol string) {

	SetPagingDefaultOptions(in)

	offset, limit, symbol = GetPagingModeByPrimaryOptions(in)

	return
}

// 设置 : 默认每页 10 条，页码 第 1 页
func SetPagingDefaultOptions(in *PageOptions) *PageOptions {

	// set default pageSize ：
	if in.PageSize < 1 {
		in.PageSize = int32(PagingSize)
	}

	// set default first page : 1
	if in.PageNumber < 1 {
		in.PageNumber = 1
	}
	// set default first currPage : 1
	if in.CurrPageNumber < 1 {
		in.CurrPageNumber = 0
	}
	// 设置默认查询字段、排序
	in.SortField, in.SortFieldTo = SetPagingModeByPrimarySelectFieldAndSort(in.SortField, in.SortFieldTo)

	return in
}

// 默认排序
func SetPagingModeByPrimarySelectFieldAndSort(SortField, SortFieldTo string) (field string, sort string) {
	sort = "desc"
	// 排序字段
	field = strings.Trim(SortField, " ")
	if field == "" {
		field = "id"
	}
	if field != "id" {
		// snake string
		field = utils.SnakeString(field)
		// 是否有排序后缀 _sort
		if field != "id" && ! strings.HasSuffix(field, SortFieldSuffix) {
			field = field + SortFieldSuffix
		}
	}
	// 排序方式
	SortFieldTo = strings.ToLower(strings.Trim(SortFieldTo, " "))
	if SortFieldTo == "asc" {
		sort = SortFieldTo
	}

	return field, sort
}

// structPointer 必须是 struct 的 指针
func PagingOptionsFieldNameIsValid(structPointer interface{}, in *PageOptions) bool {
	sortField, _ := SetPagingModeByPrimarySelectFieldAndSort(in.SortField, in.SortFieldTo)

	sElem := reflect.ValueOf(structPointer).Elem()

	return sElem.FieldByName(utils.CamelString(sortField)).IsValid()
}

/*
 * select * from users where id > ? order by id asc limit 0,PageSize;
 * select * from users where id < ? order by id desc limit 0,PageSize;
 *
 * offset,limit,symbol =  GetPagingOptions(in *PageOptions)
 * sortField := SellPointLimitTable + "." + requestParams["SortField"].(string)
 * sortFieldTo := requestParams["SortFieldTo"].(string)
 * orderBy = sortField + " " + sortFieldTo
 * sql = sql.Where(sortField + symbol + " ?", requestParams["SortValue"]).Order(orderBy)
        .Offset(offset).Limit(limit).
 */
// 主键分页模式与数字分页模式结合
func GetPagingModeByPrimaryOptions(in *PageOptions) (offset, limit int32,symbol string) {
	offset = 0
	limit = in.PageSize
	symbol = " < "
	if in.SortFieldTo == "desc" {
		if in.CurrPageNumber == 0 {
			symbol = " < "
			if in.SortValue == 0 {
				offset = (in.PageNumber -1 ) * limit
			}
		}else if in.PageNumber > in.CurrPageNumber { // 向下翻页
			offset = (in.PageNumber - in.CurrPageNumber -1 ) * limit
			symbol = " < "
		} else { //向上翻页
			offset = (in.PageNumber - 1) * limit
			symbol = " >= "
		}
	} else {
		if in.CurrPageNumber == 0 {
			symbol = " > "
			if in.SortValue == 0 {
				offset = (in.PageNumber -1 ) * limit
			}
		}else if in.PageNumber > in.CurrPageNumber { // 向下翻页
			offset = (in.PageNumber - in.CurrPageNumber -1 ) * limit
			symbol = " > "
		} else { //向上翻页
			offset = (in.PageNumber - 1) * limit
			symbol = " <= "
		}
	}


	return
}

//panic if s is not a struct pointer
func GetSortValue(s interface{}, in *PageOptions) int64 {
	SetPagingDefaultOptions(in)
	sortField := utils.CamelString(in.SortField)
	sElem := reflect.ValueOf(s).Elem()
	// 是否存在字段
	if ! sElem.FieldByName(sortField).IsValid() {
		return 0
	}
	// 字段值
	value := sElem.FieldByName(sortField).Interface()
	// 判断字段值
	switch value.(type) {
	case int64:
		return value.(int64)
	case int32:
		return int64(value.(int32))
	case int:
		return int64(value.(int))
	default:
		return 0
	}
	return 0
}

// Set Paging Result
func SetPagingResult(in *PageOptions, TotalRecords int32, SortValue int64) (paginate PageResult) {
	SetPagingDefaultOptions(in)

	paginate.TotalRecords = TotalRecords

	if paginate.TotalRecords%in.PageSize == 0 {
		paginate.TotalPages = paginate.TotalRecords / in.PageSize
	} else {
		paginate.TotalPages = paginate.TotalRecords/in.PageSize + 1
	}

	paginate.PageSize = in.PageSize
	paginate.PageNumber = in.PageNumber
	paginate.SortValue = SortValue

	return
}