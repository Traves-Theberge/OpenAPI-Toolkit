package ui

import (
    "testing"

    "github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

func TestInitialModelsAndViews(t *testing.T) {
    v := InitialValidateModel()
    // ValidateModel should have a non-empty placeholder on its TextInput
    if v.TextInput.Placeholder == "" {
        t.Fatal("InitialValidateModel TextInput placeholder is empty")
    }

    tm := InitialTestModel()
    // TestModel should have spec and url placeholders set
    if tm.SpecInput.Placeholder == "" {
        t.Fatal("InitialTestModel SpecInput placeholder is empty")
    }
    if tm.UrlInput.Placeholder == "" {
        t.Fatal("InitialTestModel UrlInput placeholder is empty")
    }

    crm := InitialCustomRequestModel()
    if crm.MethodInput.Placeholder == "" {
        t.Fatal("InitialCustomRequestModel MethodInput placeholder is empty")
    }
    if crm.EndpointInput.Placeholder == "" {
        t.Fatal("InitialCustomRequestModel EndpointInput placeholder is empty")
    }

    esm := InitialEndpointSelectorModel()
    if esm.SearchInput.Placeholder == "" {
        t.Fatal("InitialEndpointSelectorModel SearchInput placeholder is empty")
    }
    if esm.Ready {
        t.Fatal("InitialEndpointSelectorModel Ready should be false by default")
    }

    cem := InitialConfigEditorModel(models.Config{})
    if cem.SpecPathInput.Placeholder == "" {
        t.Fatal("InitialConfigEditorModel SpecPathInput placeholder is empty")
    }

    // Call view functions to ensure they render without panic
    // We only assert they return non-empty strings
    // Create a minimal model with History initialized to avoid nil pointer
    testModel := models.Model{
        History: &models.TestHistory{
            Entries: []models.HistoryEntry{},
        },
    }
    
    if ViewMenu(testModel) == "" {
        t.Fatal("ViewMenu returned empty string")
    }
    if ViewHelp(testModel) == "" {
        t.Fatal("ViewHelp returned empty string")
    }
    if ViewValidate(testModel) == "" {
        t.Fatal("ViewValidate returned empty string")
    }
    if ViewTest(testModel) == "" {
        t.Fatal("ViewTest returned empty string")
    }
    if ViewHistory(testModel) == "" {
        t.Fatal("ViewHistory returned empty string")
    }
}
