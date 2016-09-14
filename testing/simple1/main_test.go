package myPackage

import "testing"

func TestCount(t *testing.T) {
    if Count() != 1 {
        t.Error("Exepcted 1")
    } else {t.Log("Count() == 1")}
    if count != 1 {
        t.Error("Expected 1 for count too")
    } else {t.Log("count == 1")}
    if Count() != 2 {
        t.Error("Exepcted 2")
    } else {t.Log("Count() == 2")}
    if count != 2 {
        t.Error("Expected 2 for count too")
    } else {t.Log("count == 2")}
}
