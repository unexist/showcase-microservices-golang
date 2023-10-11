Feature: Create a todo
  Create various todo entries to test the endpoint.

  Scenario Outline: Create a todo with title and description and check the id.
    Given I create a todo
    When its title is "<title>"
    And its description is "<description>"
    Then its id should be <id>

    Examples:
      | title  | description  | id |
      | title1 | description1 | 1  |
      | title2 | description2 | 2  |