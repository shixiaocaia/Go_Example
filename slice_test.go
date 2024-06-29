package slice

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func printLengthAndCapacity(s []int) {
	fmt.Println(s)
	fmt.Printf("addr = %p, len = %d, cap = %d\n", &s, len(s), cap(s))
}

func TestSliceAppend(t *testing.T) {
	doAppend := func(s []int) {
		s = append(s, 1)
		fmt.Printf("addr %p\n", &s)
		printLengthAndCapacity(s)
	}

	// len = 8,所以默认分配了8个0
	s := make([]int, 8, 8)
	doAppend(s[:4])
	// 修改了原s
	assert.Equal(t, 8, len(s))
	assert.Equal(t, 1, s[4])

	doAppend(s)
	assert.NotEqual(t, 9, len(s), "大于cap扩容, 分配一个新的slice, 不会影响原slice")

	slice := make([]int, 0, 1)
	slice = append(slice, 1)
	assert.Equal(t, 1, len(slice))

	slice2 := make([]int, 1, 1)
	slice2 = append(slice2, 1)
	assert.Equal(t, 2, len(slice2), "超过cap, 会扩容")

	slice3 := make([]int, len(slice2))
	copy(slice3, slice2)
	// slice 不能直接比较
	if !reflect.DeepEqual(slice2, slice3) {
		t.Errorf("slice2 addr = %p, slice3 addr = %p", &slice2, &slice3)
	}

}

func TestSliceAppend2(t *testing.T) {
	a := make([]int, 0, 3)

	b := append(a, 1)
	printLengthAndCapacity(a)
	printLengthAndCapacity(b)
	assert.Equal(t, 1, b[0])
	assert.Equal(t, 0, len(a))

	x := append(b, 3)
	// a _ _ _
	// b 1 _ _
	// x 1 3 _
	assert.Equal(t, 1, b[0])
	assert.Equal(t, 1, x[0])

	x[0] = 10
	// b 10 _ _
	// x 10 3 _
	assert.NotEqual(t, 1, b[0], "基于同一片地址, 修改x会影响b")
	assert.Equal(t, 10, x[0])

	b = append(b, 12)
	// b 10 12 _
	// x 10 12 _
	assert.Equal(t, 12, b[1])
	assert.NotEqual(t, 3, x[1], "b append 12, x受到影响")

	b = append(b, 11, 27)
	// b 10 12 11 27
	// x 10 12 _
	assert.Equal(t, 4, len(b))
	assert.NotEqual(t, 4, len(x), "b 发生扩容")

}

func TestSliceAppend3(t *testing.T) {

	slice1 := make([]*int, 0)
	for i := 0; i < 5; i++ {
		slice1 = append(slice1, &i)
	}

	num := make([]int, 0)
	for _, v := range slice1 {
		num = append(num, *v)
	}
	// Go1.21之前是5，5，5，5，5
	assert.NotEqual(t, []int{5, 5, 5, 5, 5}, num)
}
