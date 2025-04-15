package datafetcher

type BaseUrlParams struct {
	Host     string
	User     string
	Password string
}

type PmodeUrlParams struct {
	BaseUrlParams
	Machine string
	Point   string
	Pmode   string
}

// NewPmodeUrlParams creates a new instance of PmodeUrlParams with the provided
// host, machine, point, pmode, user, and password values. It initializes the
// BaseUrlParams with the host, user, and password, and sets the machine, point,
// and pmode fields specific to PmodeUrlParams.
//
// Parameters:
//   - host: The host address.
//   - machine: The machine identifier.
//   - point: The point identifier.
//   - pmode: The pmode value.
//   - user: The username for authentication.
//   - password: The password for authentication.
//
// Returns:
//
//	A PmodeUrlParams struct populated with the provided values.
func NewPmodeUrlParams(host, machine, point, pmode, user, password string) PmodeUrlParams {
	return PmodeUrlParams{
		BaseUrlParams: BaseUrlParams{
			Host:     host,
			User:     user,
			Password: password,
		},
		Machine: machine,
		Point:   point,
		Pmode:   pmode,
	}
}

type PmodeUrlTimeParams struct {
	PmodeUrlParams
	DateTime string
}

// NewPmodeUrlTimeParams creates a new instance of PmodeUrlTimeParams with the provided
// parameters. It initializes the PmodeUrlParams using the NewPmodeUrlParams function
// and sets the DateTime field to the specified time.
//
// Parameters:
//   - host: The host address.
//   - machine: The machine identifier.
//   - point: The point identifier.
//   - pmode: The mode of operation.
//   - time: The time value to be set in the DateTime field.
//   - user: The username for authentication.
//   - password: The password for authentication.
//
// Returns:
//
//	A PmodeUrlTimeParams struct populated with the provided values.
func NewPmodeUrlTimeParams(
	host, machine, point, pmode, time, user, password string,
) PmodeUrlTimeParams {
	return PmodeUrlTimeParams{
		PmodeUrlParams: NewPmodeUrlParams(host, machine, point, pmode, user, password),
		DateTime:       time,
	}
}
