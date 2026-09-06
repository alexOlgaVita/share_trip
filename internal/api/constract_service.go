package api

type ContractRejectReason string

const (
	ReasonNoActiveContract     ContractRejectReason = "NO_ACTIVE_CONTRACT"
	ReasonContractSuspended    ContractRejectReason = "CONTRACT_SUSPENDED"
	ReasonContractTerminated   ContractRejectReason = "CONTRACT_TERMINATED"
	ReasonContractNotStarted   ContractRejectReason = "CONTRACT_NOT_STARTED"
	ReasonContractExpired      ContractRejectReason = "SERVICE_NOT_IN_CONTRACT"
	ReasonServiceNotInContract ContractRejectReason = "CONTRACT_EXPIRED"
	ReasonServiceDisabled      ContractRejectReason = "SERVICE_DISABLED"
)
