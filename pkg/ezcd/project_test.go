package ezcd_test

import (
	"testing"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
)

func TestGetProject(t *testing.T) {
	mockDB := ezcd.NewMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	mockDB.Projects["test-id"] = ezcd.Project{
		ID: "test-id",
	}

	projectID := "test-id"
	project, err := service.GetProject(projectID)


	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if project == nil {
		t.Fatalf("expected project, got nil")
	}
	if project.ID != projectID {
		t.Errorf("expected project ID %v, got %v", projectID, project.ID)
	}
}
func TestGetProjects(t *testing.T) {
	mockDB := ezcd.NewMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	mockDB.Projects["test-id-1"] = ezcd.Project{
		ID: "test-id-1",
	}
	mockDB.Projects["test-id-2"] = ezcd.Project{
		ID: "test-id-2",
	}

	projects, err := service.GetProjects()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if projects == nil {
		t.Fatalf("expected projects, got nil")
	}
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %v", len(projects))
	}
}
func TestCreateProject(t *testing.T) {
	mockDB := ezcd.NewMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	projectName := "New Project"
	project, err := service.CreateProject(projectName)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if project == nil {
		t.Fatalf("expected project, got nil")
	}
	if project.Name != projectName {
		t.Errorf("expected project name %v, got %v", projectName, project.Name)
	}
	if project.ID == "" {
		t.Errorf("expected project ID to be set, got empty string")
	}
}

