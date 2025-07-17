package main

import (
	"testing"
	"time"
)

func TestCustomCache_GetProfileShouldGetProfileWhenNotExpired(t *testing.T) {
	t.Parallel()
	testProfile := &Profile{}
	cc := NewCustomCache()
	cc.SetUpdateProfile("test1", testProfile)
	profile := cc.GetProfile("test1")
	if profile != testProfile {
		t.Error("Expected profile to be returned")
	}
}

func TestCustomCache_GetProfileShouldReturnNilWhenExpired(t *testing.T) {
	t.Parallel()

	testProfile := &Profile{}

	cc := NewCustomCache()
	cc.SetUpdateProfile("test2", testProfile)
	time.Sleep(TtlSec * time.Second)
	profile := cc.GetProfile("test2")
	if profile != nil {
		t.Error("Expected nil to be returned")
	}
}

func TestCustomCache_GetProfileShouldReturnNilIfConversionIsNotCorrect(t *testing.T) {
	t.Parallel()
	cc := NewCustomCache()
	cc.syncedMap.Store("test3", "test")
	profile := cc.GetProfile("test3")
	if profile != nil {
		t.Error("Expected nil to be returned")
	}
}

func TestCustomCache_GetProfileShouldReturnNilWhenProfileNotExists(t *testing.T) {
	t.Parallel()
	cc := NewCustomCache()
	profile := cc.GetProfile("test4")
	if profile != nil {
		t.Error("Expected nil to be returned")
	}
}

func TestCustomCache_SetProfileShouldSetProfile(t *testing.T) {
	t.Parallel()
	testProfile := &Profile{}
	cc := NewCustomCache()
	cc.SetUpdateProfile("test5", testProfile)
	profile := cc.GetProfile("test5")
	if profile != testProfile {
		t.Error("Expected profile to be returned")
	}
}

func TestCustomCache_DeleteProfileShouldDeleteProfile(t *testing.T) {
	t.Parallel()
	testProfile := &Profile{}
	cc := NewCustomCache()
	cc.SetUpdateProfile("test6", testProfile)
	cc.DeleteProfile("test6")
	profile := cc.GetProfile("test6")
	if profile != nil {
		t.Error("Expected nil to be returned")
	}
}

func TestCustomCache_AutomaticCleanUpShouldCleanUpExpiredProfiles(t *testing.T) {
	t.Parallel()
	testProfile := &Profile{}
	cc := NewCustomCache()
	cc.SetUpdateProfile("test7", testProfile)
	time.Sleep(TtlSec * time.Second)
	profile := cc.GetProfile("test")
	if profile != nil {
		t.Error("Expected nil to be returned")
	}
}
