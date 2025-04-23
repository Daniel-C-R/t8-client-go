package datafetcher

// PmodeIdentifier es una estructura que identifica un punto de medición basado en
// máquina, punto y modo de procesamiento.
type PmodeIdentifier struct {
	Machine string
	Point   string
	Pmode   string
}

// NewPmodeIdentifier crea una nueva instancia de PmodeIdentifier con los valores proporcionados
func NewPmodeIdentifier(machine, point, pmode string) PmodeIdentifier {
	return PmodeIdentifier{
		Machine: machine,
		Point:   point,
		Pmode:   pmode,
	}
}

// PmodeTimeIdentifier es una estructura que extiende PmodeIdentifier añadiendo
// información de fecha y hora para identificar un registro específico.
type PmodeTimeIdentifier struct {
	PmodeIdentifier
	DateTime string
}

func NewPmodeTimeIdentifier(machine, point, pmode, time string) PmodeTimeIdentifier {
	return PmodeTimeIdentifier{
		PmodeIdentifier: PmodeIdentifier{
			Machine: machine,
			Point:   point,
			Pmode:   pmode,
		},
		DateTime: time,
	}
}

// BaseUrlParams contiene los parámetros base para conectarse a un servidor HTTP
type BaseUrlParams struct {
	Host     string
	User     string
	Password string
}
