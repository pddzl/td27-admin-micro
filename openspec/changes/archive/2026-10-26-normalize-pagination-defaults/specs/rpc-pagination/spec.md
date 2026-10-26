## ADDED Requirements

### Requirement: Pagination defaults are applied centrally

Every list RPC in `rpc/basis` MUST obtain its effective pagination via
`common.NewPageInfo(page, pageSize)`, which clamps `Page <= 0` to 1 and
`PageSize <= 0` to 10. List logics MUST NOT construct `common.PageInfo`
directly or re-implement these defaults inline.

#### Scenario: Missing pagination values fall back to defaults

- **WHEN** a list RPC is called with page or pageSize absent or <= 0
- **THEN** the query executes with page 1 and pageSize 10
- **AND** the response contains both a populated `list` and the unfiltered
  `total`

#### Scenario: Explicit pagination values are preserved

- **WHEN** a list RPC is called with page 2 and pageSize 20
- **THEN** the query executes with offset (2-1)*20 and limit 20
