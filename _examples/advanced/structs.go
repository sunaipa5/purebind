package advanced

// SystemStatus
const (
	STATUS_OK = 0
	STATUS_ERROR = -1
	STATUS_PENDING = 0
	STATUS_UNKNOWN = 99
)

type Counter = int
type UltraCounter = Counter
type Texture2D = Texture
type TextureCubemap = Texture

type KeywordTest struct {
	TypeField int
	FuncField float32
	RangeField *byte
	SelectField int
}

type ArrayTest struct {
	Transform [16]float32
	Shadow_map [1024]uint8
	Tags [8][32]byte
}

type CallbackTest struct {
	On_event uintptr
	Calculate uintptr
}

type Texture struct {
	Id uint
	Width int
	Height int
	Mipmaps int
	Format int
}

type PaddingTest struct {
	A byte
	B int
	C byte
	D float64
}

