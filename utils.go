package actors

import (
    "reflect"
)

func toActor(p interface{}) *Actor {
    val := reflect.ValueOf(p)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    typ := val.Type()
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i).Type

        if fieldType == reflect.TypeOf(Actor{}) && field.CanAddr() && field.Addr().CanInterface() {
            return field.Addr().Interface().(*Actor)
        }
    }

    return nil
}
