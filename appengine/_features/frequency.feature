Feature:
    In order to get indicators of "Deployment Frequency"
    As a DevOps practitioner

    Scenario: should fetch actions runs
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/runs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/runs",
                "next_token": "@number@"
            }
            """

    Scenario: should get deployment frequency via runs
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._workflow_runs` WHERE date = "2021-12-30" LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | workflow_id | workflow_name | date       | runs | duration_avg       |
            | itsubaki | ghz        | 16163576    | tests         | 2021-12-30 | 10   | 0.7999999999999999 |

    Scenario: should fetch actions jobs
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/jobs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/jobs",
                "next_token": "@number@"
            }
            """

    Scenario: should get deployment frequency via jobs
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._workflow_jobs` WHERE date = "2021-12-30" LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | workflow_id | workflow_name | job_name             | date       | runs | duration_avg       |
            | itsubaki | ghz        | 16163576    | tests         | test (ubuntu-latest) | 2021-12-30 | 6    | 0.6666666666666667 |

    Scenario: should fetch releases
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/releases"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/releases",
                "next_token": "@number@"
            }
            """
