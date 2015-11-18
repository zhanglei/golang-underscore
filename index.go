package underscore

import (
	"errors"
	"reflect"
)

func Index(source, indexSelector interface{}) (interface{}, error) {
	selectorRV := reflect.ValueOf(indexSelector)
	if selectorRV.Kind() != reflect.Func {
		return nil, errors.New("underscore: Index's indexSelector is not func")
	}

	var mapRV reflect.Value
	err := each(source, func (args []reflect.Value) (bool, reflect.Value) {
		if !mapRV.IsValid() {
			mapRV = makeMapRV(selectorRV.Type().Out(0), args[0].Type())
		}

		values := selectorRV.Call(args)
		if !isErrorRVValid(values[1]) {
			mapRV.SetMapIndex(values[0], args[0])
		}

		return false, values[1]
	})
	if err == nil && mapRV.IsValid() {
		return mapRV.Interface(), nil
	}

	return nil, err
}

func IndexBy(source interface{}, property string) (interface{}, error) {
	var mapRV reflect.Value
	err := each(source, func (args []reflect.Value) (bool, reflect.Value) {
		pRV, err := getPropertyRV(args[0], property)
		if err == nil {
			if !mapRV.IsValid() {
				mapRV = makeMapRV(pRV.Type(), args[0].Type())
			}

			mapRV.SetMapIndex(pRV, args[0])
		}
		return false, reflect.ValueOf(err)
	})
	if err == nil && mapRV.IsValid() {
		return mapRV.Interface(), nil
	}

	return nil, err
	//return Index(source, func (item, _ interface{}) (interface{}, error) {
		//return getPropertyValue(item, property)
	//})
}

//Chain
func (this *Query) Index(indexSelector interface{}) Queryer {
	if this.err == nil {
		this.source, this.err = Index(this.source, indexSelector)
	}
	return this
}

func (this *Query) IndexBy(property string) Queryer {
	if this.err == nil {
		this.source, this.err = IndexBy(this.source, property)
	}
	return this
}