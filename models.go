package main

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	OrganizationID uint   `json:"organizationId" gorm:"not null"`
	Name           string `json:"name" gorm:"not null"`
}

type FinancialRecord struct {
	gorm.Model
	OrganizationID uint      `json:"organizationId" gorm:"not null"`
	Direction      string    `json:"direction" gorm:"not null"` // "IN" or "OUT"
	Amount         float64   `json:"amount" gorm:"not null"`
	Tags           []Tag     `json:"tags" gorm:"many2many:financial_record_tags;"`
	DueDate        time.Time `json:"dueDate" gorm:"not null"`
}

type CashFlowReport struct {
	MonthlyData []MonthlyCashFlow `json:"monthlyData"`
}

type MonthlyCashFlow struct {
	Year  int     `json:"year"`
	Month int     `json:"month"`
	In    float64 `json:"in"`
	Out   float64 `json:"out"`
}

// E-commerce Domain
type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"not null"`
	SKU         string  `json:"sku" gorm:"uniqueIndex"`
	Stock       int     `json:"stock" gorm:"not null;default:0"`
	CategoryID  uint    `json:"categoryId"`
	Images      []Image `json:"images" gorm:"many2many:product_images;"`
}

type Category struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parentId"`
	Products    []Product `json:"products" gorm:"foreignKey:CategoryID"`
}

type Customer struct {
	gorm.Model
	FirstName   string    `json:"firstName" gorm:"not null"`
	LastName    string    `json:"lastName" gorm:"not null"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber string    `json:"phoneNumber"`
	Address     []Address `json:"addresses" gorm:"foreignKey:CustomerID"`
	Orders      []Order   `json:"orders" gorm:"foreignKey:CustomerID"`
}

type Order struct {
	gorm.Model
	CustomerID    uint          `json:"customerId" gorm:"not null"`
	OrderDate     time.Time     `json:"orderDate" gorm:"not null"`
	Status        string        `json:"status" gorm:"not null;default:'pending'"`
	TotalAmount   float64       `json:"totalAmount" gorm:"not null"`
	ShippingCost  float64       `json:"shippingCost"`
	TrackingCode  string        `json:"trackingCode"`
	PaymentMethod string        `json:"paymentMethod"`
	OrderItems    []OrderItem   `json:"orderItems" gorm:"foreignKey:OrderID"`
	Address       Address       `json:"address" gorm:"embedded"`
	Transactions  []Transaction `json:"transactions" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID    uint    `json:"orderId" gorm:"not null"`
	ProductID  uint    `json:"productId" gorm:"not null"`
	Quantity   int     `json:"quantity" gorm:"not null"`
	UnitPrice  float64 `json:"unitPrice" gorm:"not null"`
	TotalPrice float64 `json:"totalPrice" gorm:"not null"`
}

// Healthcare Domain
type Patient struct {
	gorm.Model
	FirstName      string         `json:"firstName" gorm:"not null"`
	LastName       string         `json:"lastName" gorm:"not null"`
	DateOfBirth    time.Time      `json:"dateOfBirth" gorm:"not null"`
	Gender         string         `json:"gender"`
	Email          string         `json:"email" gorm:"uniqueIndex"`
	PhoneNumber    string         `json:"phoneNumber"`
	Address        Address        `json:"address" gorm:"embedded"`
	BloodType      string         `json:"bloodType"`
	Allergies      string         `json:"allergies"`
	MedicalHistory string         `json:"medicalHistory" gorm:"type:text"`
	Appointments   []Appointment  `json:"appointments" gorm:"foreignKey:PatientID"`
	Prescriptions  []Prescription `json:"prescriptions" gorm:"foreignKey:PatientID"`
}

type Doctor struct {
	gorm.Model
	FirstName    string        `json:"firstName" gorm:"not null"`
	LastName     string        `json:"lastName" gorm:"not null"`
	Specialty    string        `json:"specialty" gorm:"not null"`
	Email        string        `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber  string        `json:"phoneNumber"`
	LicenseID    string        `json:"licenseId" gorm:"uniqueIndex;not null"`
	Department   string        `json:"department"`
	Appointments []Appointment `json:"appointments" gorm:"foreignKey:DoctorID"`
}

type Appointment struct {
	gorm.Model
	PatientID  uint      `json:"patientId" gorm:"not null"`
	DoctorID   uint      `json:"doctorId" gorm:"not null"`
	StartTime  time.Time `json:"startTime" gorm:"not null"`
	EndTime    time.Time `json:"endTime" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null;default:'scheduled'"`
	Notes      string    `json:"notes" gorm:"type:text"`
	RoomNumber string    `json:"roomNumber"`
}

type Prescription struct {
	gorm.Model
	PatientID      uint      `json:"patientId" gorm:"not null"`
	DoctorID       uint      `json:"doctorId" gorm:"not null"`
	IssuedDate     time.Time `json:"issuedDate" gorm:"not null"`
	ExpirationDate time.Time `json:"expirationDate" gorm:"not null"`
	Instructions   string    `json:"instructions" gorm:"type:text;not null"`
	Medications    []Medication `json:"medications" gorm:"many2many:prescription_medications;"`
}

type Medication struct {
	gorm.Model
	Name          string `json:"name" gorm:"not null"`
	Description   string `json:"description" gorm:"type:text"`
	Dosage        string `json:"dosage" gorm:"not null"`
	Manufacturer  string `json:"manufacturer"`
	SideEffects   string `json:"sideEffects" gorm:"type:text"`
	Contraindications string `json:"contraindications" gorm:"type:text"`
}

// Education Domain
type Student struct {
	gorm.Model
	FirstName      string      `json:"firstName" gorm:"not null"`
	LastName       string      `json:"lastName" gorm:"not null"`
	Email          string      `json:"email" gorm:"uniqueIndex;not null"`
	DateOfBirth    time.Time   `json:"dateOfBirth"`
	EnrollmentDate time.Time   `json:"enrollmentDate" gorm:"not null"`
	GraduationDate *time.Time  `json:"graduationDate"`
	Major          string      `json:"major"`
	Address        Address     `json:"address" gorm:"embedded"`
	Courses        []Course    `json:"courses" gorm:"many2many:student_courses;"`
	Grades         []Grade     `json:"grades" gorm:"foreignKey:StudentID"`
}

type Course struct {
	gorm.Model
	Code        string     `json:"code" gorm:"uniqueIndex;not null"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text"`
	Credits     int        `json:"credits" gorm:"not null"`
	Department  string     `json:"department" gorm:"not null"`
	TeacherID   uint       `json:"teacherId"`
	Students    []Student  `json:"students" gorm:"many2many:student_courses;"`
	Assignments []Assignment `json:"assignments" gorm:"foreignKey:CourseID"`
}

type Teacher struct {
	gorm.Model
	FirstName   string   `json:"firstName" gorm:"not null"`
	LastName    string   `json:"lastName" gorm:"not null"`
	Email       string   `json:"email" gorm:"uniqueIndex;not null"`
	Department  string   `json:"department" gorm:"not null"`
	HireDate    time.Time `json:"hireDate" gorm:"not null"`
	Credentials string   `json:"credentials"`
	Courses     []Course `json:"courses" gorm:"foreignKey:TeacherID"`
}

type Grade struct {
	gorm.Model
	StudentID    uint    `json:"studentId" gorm:"not null"`
	CourseID     uint    `json:"courseId" gorm:"not null"`
	AssignmentID *uint   `json:"assignmentId"`
	Score        float64 `json:"score" gorm:"not null"`
	Comments     string  `json:"comments"`
	GradedDate   time.Time `json:"gradedDate" gorm:"not null"`
}

// Shared model across domains
type Address struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	CustomerID   *uint  `json:"customerId"`
	Street       string `json:"street" gorm:"not null"`
	City         string `json:"city" gorm:"not null"`
	State        string `json:"state" gorm:"not null"`
	ZipCode      string `json:"zipCode" gorm:"not null"`
	Country      string `json:"country" gorm:"not null"`
	IsDefault    bool   `json:"isDefault" gorm:"default:false"`
	AddressType  string `json:"addressType" gorm:"default:'shipping'"`
}

type Image struct {
	gorm.Model
	URL         string `json:"url" gorm:"not null"`
	Description string `json:"description"`
	AltText     string `json:"altText"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type Transaction struct {
	gorm.Model
	OrderID       *uint     `json:"orderId"`
	Amount        float64   `json:"amount" gorm:"not null"`
	Currency      string    `json:"currency" gorm:"not null;default:'USD'"`
	PaymentMethod string    `json:"paymentMethod" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null"`
	ProcessedAt   time.Time `json:"processedAt" gorm:"not null"`
	ReferenceCode string    `json:"referenceCode"`
}

// Transportation Domain
type Vehicle struct {
	gorm.Model
	Make        string `json:"make" gorm:"not null"`
	Year        int    `json:"year" gorm:"not null"`
	LicensePlate string `json:"licensePlate" gorm:"uniqueIndex"`
	Color       string `json:"color"`
	VIN         string `json:"vin" gorm:"uniqueIndex;not null"`
	OwnerID     uint   `json:"ownerId"`
	Trips       []Trip `json:"trips" gorm:"foreignKey:VehicleID"`
}

type Driver struct {
	gorm.Model
	FirstName      string    `json:"firstName" gorm:"not null"`
	LastName       string    `json:"lastName" gorm:"not null"`
	LicenseNumber  string    `json:"licenseNumber" gorm:"uniqueIndex;not null"`
	LicenseExpiry  time.Time `json:"licenseExpiry" gorm:"not null"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phoneNumber" gorm:"not null"`
	DateOfBirth    time.Time `json:"dateOfBirth" gorm:"not null"`
	Rating         float64   `json:"rating"`
	Trips          []Trip    `json:"trips" gorm:"foreignKey:DriverID"`
}

type Trip struct {
	gorm.Model
	VehicleID       uint      `json:"vehicleId" gorm:"not null"`
	DriverID        uint      `json:"driverId"`
	StartLocation   string    `json:"startLocation" gorm:"not null"`
	EndLocation     string    `json:"endLocation" gorm:"not null"`
	StartTime       time.Time `json:"startTime" gorm:"not null"`
	EndTime         *time.Time `json:"endTime"`
	Distance        float64   `json:"distance"`
	Fare            float64   `json:"fare"`
	Status          string    `json:"status" gorm:"not null;default:'scheduled'"`
	PassengerID     *uint     `json:"passengerId"`
}

type Passenger struct {
	gorm.Model
	FirstName   string `json:"firstName" gorm:"not null"`
	LastName    string `json:"lastName" gorm:"not null"`
	Email       string `json:"email" gorm:"uniqueIndex"`
	PhoneNumber string `json:"phoneNumber" gorm:"not null"`
	Rating      float64 `json:"rating"`
	Trips       []Trip `json:"trips" gorm:"foreignKey:PassengerID"`
}

// Social Media Domain
type User struct {
	gorm.Model
	Username      string    `json:"username" gorm:"uniqueIndex;not null"`
	Email         string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash  string    `json:"passwordHash" gorm:"not null"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Bio           string    `json:"bio" gorm:"type:text"`
	ProfilePicURL string    `json:"profilePicUrl"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	JoinDate      time.Time `json:"joinDate" gorm:"not null"`
	LastLogin     time.Time `json:"lastLogin"`
	Posts         []Post    `json:"posts" gorm:"foreignKey:UserID"`
	Followers     []Follow  `json:"followers" gorm:"foreignKey:FolloweeID"`
	Following     []Follow  `json:"following" gorm:"foreignKey:FollowerID"`
}

type Post struct {
	gorm.Model
	UserID      uint      `json:"userId" gorm:"not null"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	PostDate    time.Time `json:"postDate" gorm:"not null"`
	Likes       int       `json:"likes" gorm:"default:0"`
	Comments    []Comment `json:"comments" gorm:"foreignKey:PostID"`
	Attachments []Attachment `json:"attachments" gorm:"foreignKey:PostID"`
	IsPublic    bool      `json:"isPublic" gorm:"default:true"`
}

type Comment struct {
	gorm.Model
	PostID      uint      `json:"postId" gorm:"not null"`
	UserID      uint      `json:"userId" gorm:"not null"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	CommentDate time.Time `json:"commentDate" gorm:"not null"`
	Likes       int       `json:"likes" gorm:"default:0"`
	ParentID    *uint     `json:"parentId"`
}

type Follow struct {
	gorm.Model
	FollowerID  uint      `json:"followerId" gorm:"not null"`
	FolloweeID  uint      `json:"followeeId" gorm:"not null"`
	FollowDate  time.Time `json:"followDate" gorm:"not null"`
}

type Attachment struct {
	gorm.Model
	PostID      uint   `json:"postId" gorm:"not null"`
	URL         string `json:"url" gorm:"not null"`
	Type        string `json:"type" gorm:"not null"` // image, video, audio, etc.
	Description string `json:"description"`
}

// Human Resources Domain
type Employee struct {
	gorm.Model
	FirstName       string      `json:"firstName" gorm:"not null"`
	LastName        string      `json:"lastName" gorm:"not null"`
	Email           string      `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber     string      `json:"phoneNumber"`
	HireDate        time.Time   `json:"hireDate" gorm:"not null"`
	TerminationDate *time.Time  `json:"terminationDate"`
	PositionID      uint        `json:"positionId" gorm:"not null"`
	DepartmentID    uint        `json:"departmentId" gorm:"not null"`
	ManagerID       *uint       `json:"managerId"`
	Salary          float64     `json:"salary" gorm:"not null"`
	Address         Address     `json:"address" gorm:"embedded"`
	Leaves          []Leave     `json:"leaves" gorm:"foreignKey:EmployeeID"`
	Evaluations     []Evaluation `json:"evaluations" gorm:"foreignKey:EmployeeID"`
}

type Department struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	ManagerID   uint       `json:"managerId"`
	Employees   []Employee `json:"employees" gorm:"foreignKey:DepartmentID"`
	Budget      float64    `json:"budget"`
}

type Position struct {
	gorm.Model
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text"`
	MinSalary   float64    `json:"minSalary"`
	MaxSalary   float64    `json:"maxSalary"`
	Employees   []Employee `json:"employees" gorm:"foreignKey:PositionID"`
}

type Leave struct {
	gorm.Model
	EmployeeID  uint      `json:"employeeId" gorm:"not null"`
	StartDate   time.Time `json:"startDate" gorm:"not null"`
	EndDate     time.Time `json:"endDate" gorm:"not null"`
	LeaveType   string    `json:"leaveType" gorm:"not null"` // sick, vacation, personal
	Status      string    `json:"status" gorm:"not null;default:'pending'"`
	Reason      string    `json:"reason"`
	ApprovedBy  *uint     `json:"approvedBy"`
}

type Evaluation struct {
	gorm.Model
	EmployeeID    uint      `json:"employeeId" gorm:"not null"`
	EvaluatorID   uint      `json:"evaluatorId" gorm:"not null"`
	EvaluationDate time.Time `json:"evaluationDate" gorm:"not null"`
	Performance   int       `json:"performance" gorm:"not null"` // 1-5 scale
	Comments      string    `json:"comments" gorm:"type:text"`
	GoalsSet      string    `json:"goalsSet" gorm:"type:text"`
}

// Content Management Domain
type Website struct {
	gorm.Model
	Name         string    `json:"name" gorm:"not null"`
	Domain       string    `json:"domain" gorm:"uniqueIndex;not null"`
	Description  string    `json:"description"`
	LaunchDate   time.Time `json:"launchDate"`
	OwnerID      uint      `json:"ownerId" gorm:"not null"`
	Pages        []Page    `json:"pages" gorm:"foreignKey:WebsiteID"`
	Theme        string    `json:"theme"`
	Settings     JSONData  `json:"settings" gorm:"type:jsonb"`
}

type Page struct {
	gorm.Model
	WebsiteID    uint      `json:"websiteId" gorm:"not null"`
	Title        string    `json:"title" gorm:"not null"`
	Slug         string    `json:"slug" gorm:"not null"`
	Content      string    `json:"content" gorm:"type:text"`
	Published    bool      `json:"published" gorm:"default:false"`
	PublishedAt  *time.Time `json:"publishedAt"`
	AuthorID     uint      `json:"authorId"`
	SEOMetadata  JSONData  `json:"seoMetadata" gorm:"type:jsonb"`
	ParentID     *uint     `json:"parentId"`
	SortOrder    int       `json:"sortOrder" gorm:"default:0"`
}

type MenuItem struct {
	gorm.Model
	WebsiteID  uint   `json:"websiteId" gorm:"not null"`
	Label      string `json:"label" gorm:"not null"`
	URL        string `json:"url"`
	PageID     *uint  `json:"pageId"`
	ParentID   *uint  `json:"parentId"`
	SortOrder  int    `json:"sortOrder" gorm:"default:0"`
	IsExternal bool   `json:"isExternal" gorm:"default:false"`
}

// IoT & Smart Home Domain
type Device struct {
	gorm.Model
	Name         string        `json:"name" gorm:"not null"`
	DeviceType   string        `json:"deviceType" gorm:"not null"`
	SerialNumber string        `json:"serialNumber" gorm:"uniqueIndex;not null"`
	HomeID       uint          `json:"homeId" gorm:"not null"`
	RoomID       *uint         `json:"roomId"`
	Status       string        `json:"status" gorm:"not null;default:'offline'"`
	IPAddress    string        `json:"ipAddress"`
	MacAddress   string        `json:"macAddress"`
	Firmware     string        `json:"firmware"`
	Readings     []DeviceReading `json:"readings" gorm:"foreignKey:DeviceID"`
	Settings     JSONData      `json:"settings" gorm:"type:jsonb"`
}

type Home struct {
	gorm.Model
	Name      string   `json:"name" gorm:"not null"`
	Address   Address  `json:"address" gorm:"embedded"`
	OwnerID   uint     `json:"ownerId" gorm:"not null"`
	Rooms     []Room   `json:"rooms" gorm:"foreignKey:HomeID"`
	Devices   []Device `json:"devices" gorm:"foreignKey:HomeID"`
	TimeZone  string   `json:"timeZone"`
}

type Room struct {
	gorm.Model
	Name    string   `json:"name" gorm:"not null"`
	HomeID  uint     `json:"homeId" gorm:"not null"`
	Floor   int      `json:"floor" gorm:"default:1"`
	Devices []Device `json:"devices" gorm:"foreignKey:RoomID"`
}

type DeviceReading struct {
	gorm.Model
	DeviceID   uint      `json:"deviceId" gorm:"not null"`
	ReadingType string    `json:"readingType" gorm:"not null"` // temperature, humidity, power, etc.
	Value      float64   `json:"value" gorm:"not null"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp" gorm:"not null"`
}

// Utility type for JSON data
type JSONData map[string]interface{}

// Fitness & Wellness Domain
type Workout struct {
	gorm.Model
	UserID      uint           `json:"userId" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	StartTime   time.Time      `json:"startTime" gorm:"not null"`
	EndTime     time.Time      `json:"endTime" gorm:"not null"`
	CaloriesBurned float64     `json:"caloriesBurned"`
	Notes       string         `json:"notes"`
	WorkoutType string         `json:"workoutType"`
	Exercises   []WorkoutExercise `json:"exercises" gorm:"foreignKey:WorkoutID"`
}

type Exercise struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	MuscleGroup string    `json:"muscleGroup"`
	Equipment   string    `json:"equipment"`
	Difficulty  string    `json:"difficulty"`
	Instructions string   `json:"instructions" gorm:"type:text"`
}

type WorkoutExercise struct {
	gorm.Model
	WorkoutID  uint    `json:"workoutId" gorm:"not null"`
	ExerciseID uint    `json:"exerciseId" gorm:"not null"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Duration   int     `json:"duration"` // in seconds
	Notes      string  `json:"notes"`
}

type MealPlan struct {
	gorm.Model
	UserID      uint      `json:"userId" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	StartDate   time.Time `json:"startDate" gorm:"not null"`
	EndDate     time.Time `json:"endDate" gorm:"not null"`
	TotalCalories float64 `json:"totalCalories"`
	Notes       string    `json:"notes"`
	Meals       []Meal    `json:"meals" gorm:"foreignKey:MealPlanID"`
}

type Meal struct {
	gorm.Model
	MealPlanID  uint      `json:"mealPlanId" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"` // breakfast, lunch, dinner, snack
	Time        time.Time `json:"time" gorm:"not null"`
	TotalCalories float64 `json:"totalCalories"`
	Foods       []MealFood `json:"foods" gorm:"foreignKey:MealID"`
}

// Finance Domain (Extending existing models)
type Budget struct {
	gorm.Model
	OrganizationID uint      `json:"organizationId" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	StartDate      time.Time `json:"startDate" gorm:"not null"`
	EndDate        time.Time `json:"endDate" gorm:"not null"`
	TotalAmount    float64   `json:"totalAmount" gorm:"not null"`
	Categories     []BudgetCategory `json:"categories" gorm:"foreignKey:BudgetID"`
}

type BudgetCategory struct {
	gorm.Model
	BudgetID      uint    `json:"budgetId" gorm:"not null"`
	Name          string  `json:"name" gorm:"not null"`
	PlannedAmount float64 `json:"plannedAmount" gorm:"not null"`
	ActualAmount  float64 `json:"actualAmount" gorm:"default:0"`
}

type Invoice struct {
	gorm.Model
	OrganizationID  uint      `json:"organizationId" gorm:"not null"`
	CustomerID      uint      `json:"customerId" gorm:"not null"`
	InvoiceNumber   string    `json:"invoiceNumber" gorm:"not null"`
	IssueDate       time.Time `json:"issueDate" gorm:"not null"`
	DueDate         time.Time `json:"dueDate" gorm:"not null"`
	Status          string    `json:"status" gorm:"not null;default:'pending'"`
	TotalAmount     float64   `json:"totalAmount" gorm:"not null"`
	InvoiceItems    []InvoiceItem `json:"invoiceItems" gorm:"foreignKey:InvoiceID"`
}

type InvoiceItem struct {
	gorm.Model
	InvoiceID     uint    `json:"invoiceId" gorm:"not null"`
	Description   string  `json:"description" gorm:"not null"`
	Quantity      int     `json:"quantity" gorm:"not null;default:1"`
	UnitPrice     float64 `json:"unitPrice" gorm:"not null"`
	TotalPrice    float64 `json:"totalPrice" gorm:"not null"`
}

// Real Estate Domain
type Property struct {
	gorm.Model
	Address         Address   `json:"address" gorm:"embedded"`
	OwnerID         uint      `json:"ownerId"`
	PropertyType    string    `json:"propertyType" gorm:"not null"` // apartment, house, commercial
	SquareFootage   float64   `json:"squareFootage"`
	Bedrooms        int       `json:"bedrooms"`
	Bathrooms       float64   `json:"bathrooms"`
	YearBuilt       int       `json:"yearBuilt"`
	ListingPrice    float64   `json:"listingPrice"`
	IsForSale       bool      `json:"isForSale" gorm:"default:false"`
	IsForRent       bool      `json:"isForRent" gorm:"default:false"`
	Features        string    `json:"features" gorm:"type:text"`
	Images          []Image   `json:"images" gorm:"many2many:property_images;"`
	Listings        []Listing `json:"listings" gorm:"foreignKey:PropertyID"`
}

type Listing struct {
	gorm.Model
	PropertyID    uint      `json:"propertyId" gorm:"not null"`
	AgentID       uint      `json:"agentId" gorm:"not null"`
	ListingDate   time.Time `json:"listingDate" gorm:"not null"`
	ExpirationDate time.Time `json:"expirationDate"`
	ListingType   string    `json:"listingType" gorm:"not null"` // sale, rent
	Price         float64   `json:"price" gorm:"not null"`
	Description   string    `json:"description" gorm:"type:text"`
	Status        string    `json:"status" gorm:"not null;default:'active'"`
}

type Agent struct {
	gorm.Model
	FirstName    string    `json:"firstName" gorm:"not null"`
	LastName     string    `json:"lastName" gorm:"not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber  string    `json:"phoneNumber" gorm:"not null"`
	LicenseNumber string   `json:"licenseNumber" gorm:"uniqueIndex;not null"`
	Agency       string    `json:"agency"`
	Listings     []Listing `json:"listings" gorm:"foreignKey:AgentID"`
}

// Hospitality Domain
type Hotel struct {
	gorm.Model
	Name           string    `json:"name" gorm:"not null"`
	Address        Address   `json:"address" gorm:"embedded"`
	Rating         float64   `json:"rating"` // 1-5 stars
	PhoneNumber    string    `json:"phoneNumber"`
	Email          string    `json:"email"`
	Website        string    `json:"website"`
	CheckInTime    string    `json:"checkInTime"`
	CheckOutTime   string    `json:"checkOutTime"`
	TotalRooms     int       `json:"totalRooms" gorm:"not null"`
	Description    string    `json:"description" gorm:"type:text"`
	Amenities      string    `json:"amenities" gorm:"type:text"`
	Rooms          []Room    `json:"rooms" gorm:"foreignKey:HotelID"`
	Bookings       []Booking `json:"bookings" gorm:"foreignKey:HotelID"`
}

type HotelRoom struct {
	gorm.Model
	HotelID      uint    `json:"hotelId" gorm:"not null"`
	RoomNumber   string  `json:"roomNumber" gorm:"not null"`
	RoomType     string  `json:"roomType" gorm:"not null"` // single, double, suite
	Price        float64 `json:"price" gorm:"not null"`
	Capacity     int     `json:"capacity" gorm:"not null"`
	IsAvailable  bool    `json:"isAvailable" gorm:"default:true"`
	Description  string  `json:"description"`
	FloorNumber  int     `json:"floorNumber"`
	Amenities    string  `json:"amenities"`
	Bookings     []Booking `json:"bookings" gorm:"foreignKey:RoomID"`
}

type Booking struct {
	gorm.Model
	HotelID      uint      `json:"hotelId" gorm:"not null"`
	RoomID       uint      `json:"roomId" gorm:"not null"`
	GuestID      uint      `json:"guestId" gorm:"not null"`
	CheckInDate  time.Time `json:"checkInDate" gorm:"not null"`
	CheckOutDate time.Time `json:"checkOutDate" gorm:"not null"`
	TotalPrice   float64   `json:"totalPrice" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null;default:'confirmed'"`
	BookingDate  time.Time `json:"bookingDate" gorm:"not null"`
	GuestCount   int       `json:"guestCount" gorm:"not null;default:1"`
	SpecialRequests string `json:"specialRequests"`
}

type Guest struct {
	gorm.Model
	FirstName    string    `json:"firstName" gorm:"not null"`
	LastName     string    `json:"lastName" gorm:"not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber  string    `json:"phoneNumber"`
	Address      Address   `json:"address" gorm:"embedded"`
	DateOfBirth  time.Time `json:"dateOfBirth"`
	IDType       string    `json:"idType"` // passport, driver's license
	IDNumber     string    `json:"idNumber"`
	Nationality  string    `json:"nationality"`
	Bookings     []Booking `json:"bookings" gorm:"foreignKey:GuestID"`
}

// Event Management Domain
type Event struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartDate   time.Time `json:"startDate" gorm:"not null"`
	EndDate     time.Time `json:"endDate" gorm:"not null"`
	Location    string    `json:"location"`
	Address     Address   `json:"address" gorm:"embedded"`
	MaxCapacity int       `json:"maxCapacity"`
	OrganizerID uint      `json:"organizerId" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null;default:'planned'"`
	EventType   string    `json:"eventType"`
	Tickets     []Ticket  `json:"tickets" gorm:"foreignKey:EventID"`
	Sessions    []EventSession `json:"sessions" gorm:"foreignKey:EventID"`
}

type Ticket struct {
	gorm.Model
	EventID     uint    `json:"eventId" gorm:"not null"`
	AttendeeID  uint    `json:"attendeeId" gorm:"not null"`
	TicketType  string  `json:"ticketType" gorm:"not null"` // VIP, standard, etc.
	Price       float64 `json:"price" gorm:"not null"`
	Status      string  `json:"status" gorm:"not null;default:'reserved'"`
	PurchaseDate time.Time `json:"purchaseDate" gorm:"not null"`
	QRCode      string  `json:"qrCode"`
}

type EventSession struct {
	gorm.Model
	EventID     uint      `json:"eventId" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartTime   time.Time `json:"startTime" gorm:"not null"`
	EndTime     time.Time `json:"endTime" gorm:"not null"`
	Location    string    `json:"location"`
	SpeakerID   *uint     `json:"speakerId"`
	MaxCapacity int       `json:"maxCapacity"`
}

type Speaker struct {
	gorm.Model
	FirstName   string    `json:"firstName" gorm:"not null"`
	LastName    string    `json:"lastName" gorm:"not null"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber string    `json:"phoneNumber"`
	Bio         string    `json:"bio" gorm:"type:text"`
	ProfilePic  string    `json:"profilePic"`
	Sessions    []EventSession `json:"sessions" gorm:"foreignKey:SpeakerID"`
}

// Manufacturing Domain
type Component struct {
	gorm.Model
	Name          string    `json:"name" gorm:"not null"`
	SKU           string    `json:"sku" gorm:"uniqueIndex;not null"`
	Description   string    `json:"description"`
	InventoryCount int      `json:"inventoryCount" gorm:"default:0"`
	ReorderLevel  int       `json:"reorderLevel"`
	Cost          float64   `json:"cost" gorm:"not null"`
	LeadTime      int       `json:"leadTime"` // in days
	SupplierID    uint      `json:"supplierId"`
}

type Supplier struct {
	gorm.Model
	Name          string    `json:"name" gorm:"not null"`
	Email         string    `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber   string    `json:"phoneNumber"`
	ContactPerson string    `json:"contactPerson"`
	Address       Address   `json:"address" gorm:"embedded"`
	PaymentTerms  string    `json:"paymentTerms"`
	Products      []Product `json:"products" gorm:"many2many:product_suppliers;"`
	Rating        float64   `json:"rating"` // 1-5 scale
}

// Logistics Domain
type Shipment struct {
	gorm.Model
	TrackingNumber string    `json:"trackingNumber" gorm:"uniqueIndex;not null"`
	Origin         string    `json:"origin" gorm:"not null"`
	Destination    string    `json:"destination" gorm:"not null"`
	ShipDate       time.Time `json:"shipDate" gorm:"not null"`
	EstimatedDelivery time.Time `json:"estimatedDelivery"`
	ActualDelivery *time.Time `json:"actualDelivery"`
	Status         string    `json:"status" gorm:"not null;default:'processing'"`
	CarrierID      uint      `json:"carrierId" gorm:"not null"`
	Weight         float64   `json:"weight"`
	Dimensions     string    `json:"dimensions"`
	ShippingCost   float64   `json:"shippingCost"`
	Packages       []Package `json:"packages" gorm:"foreignKey:ShipmentID"`
}

type Package struct {
	gorm.Model
	ShipmentID     uint    `json:"shipmentId" gorm:"not null"`
	Weight         float64 `json:"weight" gorm:"not null"`
	Length         float64 `json:"length"`
	Width          float64 `json:"width"`
	Height         float64 `json:"height"`
	Contents       string  `json:"contents"`
	DeclaredValue  float64 `json:"declaredValue"`
}

type Carrier struct {
	gorm.Model
	Name          string    `json:"name" gorm:"not null"`
	TrackingURL   string    `json:"trackingUrl"`
	ContactNumber string    `json:"contactNumber"`
	Email         string    `json:"email"`
	ServiceLevel  string    `json:"serviceLevel"` // express, standard, economy
	Shipments     []Shipment `json:"shipments" gorm:"foreignKey:CarrierID"`
}

// Project Management Domain
type Project struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartDate   time.Time `json:"startDate" gorm:"not null"`
	EndDate     time.Time `json:"endDate"`
	Status      string    `json:"status" gorm:"not null;default:'planning'"`
	ManagerID   uint      `json:"managerId" gorm:"not null"`
	ClientID    uint      `json:"clientId"`
	Budget      float64   `json:"budget"`
	Tasks       []Task    `json:"tasks" gorm:"foreignKey:ProjectID"`
	Milestones  []Milestone `json:"milestones" gorm:"foreignKey:ProjectID"`
	Team        []TeamMember `json:"team" gorm:"many2many:project_team;"`
}

type Task struct {
	gorm.Model
	ProjectID    uint      `json:"projectId" gorm:"not null"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text"`
	StartDate    time.Time `json:"startDate"`
	DueDate      time.Time `json:"dueDate"`
	Status       string    `json:"status" gorm:"not null;default:'todo'"`
	Priority     string    `json:"priority" gorm:"not null;default:'medium'"`
	AssigneeID   *uint     `json:"assigneeId"`
	ParentID     *uint     `json:"parentId"`
	EstimatedHours float64 `json:"estimatedHours"`
	ActualHours  float64   `json:"actualHours"`
	Subtasks     []Task    `json:"subtasks" gorm:"foreignKey:ParentID"`
}

type Milestone struct {
	gorm.Model
	ProjectID    uint      `json:"projectId" gorm:"not null"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description"`
	DueDate      time.Time `json:"dueDate" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null;default:'pending'"`
	Deliverables string    `json:"deliverables"`
}

// Legal Domain
type Contract struct {
	gorm.Model
	Title            string    `json:"title" gorm:"not null"`
	Description      string    `json:"description" gorm:"type:text"`
	StartDate        time.Time `json:"startDate" gorm:"not null"`
	EndDate          time.Time `json:"endDate"`
	Status           string    `json:"status" gorm:"not null;default:'draft'"`
	PartyOneID       uint      `json:"partyOneId" gorm:"not null"`
	PartyTwoID       uint      `json:"partyTwoId" gorm:"not null"`
	Value            float64   `json:"value"`
	Clauses          []Clause  `json:"clauses" gorm:"foreignKey:ContractID"`
	SignedDate       *time.Time `json:"signedDate"`
	RenewalTerms     string    `json:"renewalTerms"`
	TerminationTerms string    `json:"terminationTerms"`
}

type Clause struct {
	gorm.Model
	ContractID  uint   `json:"contractId" gorm:"not null"`
	Title       string `json:"title" gorm:"not null"`
	Content     string `json:"content" gorm:"type:text;not null"`
	SectionNumber int  `json:"sectionNumber"`
	IsRequired  bool   `json:"isRequired" gorm:"default:true"`
}

type LegalEntity struct {
	gorm.Model
	Name          string    `json:"name" gorm:"not null"`
	Type          string    `json:"type" gorm:"not null"` // individual, corporation, partnership
	ContactName   string    `json:"contactName"`
	ContactEmail  string    `json:"contactEmail"`
	ContactPhone  string    `json:"contactPhone"`
	Address       Address   `json:"address" gorm:"embedded"`
	TaxID         string    `json:"taxId"`
	Contracts     []Contract `json:"contractsAsPartyOne" gorm:"foreignKey:PartyOneID"`
	ContractsTwo  []Contract `json:"contractsAsPartyTwo" gorm:"foreignKey:PartyTwoID"`
}

// Agriculture Domain
type Farm struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	Address     Address `json:"address" gorm:"embedded"`
	OwnerID     uint    `json:"ownerId" gorm:"not null"`
	TotalArea   float64 `json:"totalArea"` // in acres/hectares
	Description string  `json:"description"`
	FarmType    string  `json:"farmType"` // livestock, crops, mixed
	Fields      []Field `json:"fields" gorm:"foreignKey:FarmID"`
	Livestock   []Livestock `json:"livestock" gorm:"foreignKey:FarmID"`
}

type Field struct {
	gorm.Model
	FarmID      uint    `json:"farmId" gorm:"not null"`
	Name        string  `json:"name" gorm:"not null"`
	Area        float64 `json:"area" gorm:"not null"` // in acres/hectares
	SoilType    string  `json:"soilType"`
	Crop        *string `json:"crop"`
	PlantingDate *time.Time `json:"plantingDate"`
	HarvestDate *time.Time `json:"harvestDate"`
	Status      string  `json:"status" gorm:"not null;default:'fallow'"`
}

type Livestock struct {
	gorm.Model
	FarmID      uint    `json:"farmId" gorm:"not null"`
	Type        string  `json:"type" gorm:"not null"` // cattle, sheep, poultry
	Breed       string  `json:"breed"`
	Count       int     `json:"count" gorm:"not null"`
	Acquisition time.Time `json:"acquisition" gorm:"not null"`
	Notes       string  `json:"notes"`
}

// Assignment (Education)
type Assignment struct {
	gorm.Model
	CourseID     uint      `json:"courseId" gorm:"not null"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text"`
	DueDate      time.Time `json:"dueDate" gorm:"not null"`
	TotalPoints  float64   `json:"totalPoints" gorm:"not null"`
	AssignmentType string  `json:"assignmentType"` // quiz, exam, project, paper
	Instructions string    `json:"instructions" gorm:"type:text"`
}

// MealFood (Fitness & Wellness)
type MealFood struct {
	gorm.Model
	MealID       uint    `json:"mealId" gorm:"not null"`
	FoodName     string  `json:"foodName" gorm:"not null"`
	Portion      float64 `json:"portion" gorm:"not null"`
	Unit         string  `json:"unit" gorm:"not null"` // grams, ounces, cups
	Calories     float64 `json:"calories"`
	Protein      float64 `json:"protein"` // in grams
	Carbs        float64 `json:"carbs"`   // in grams
	Fat          float64 `json:"fat"`     // in grams
}

// TeamMember (Project Management)
type TeamMember struct {
	gorm.Model
	UserID       uint    `json:"userId" gorm:"not null"`
	Role         string  `json:"role" gorm:"not null"`
	JoinDate     time.Time `json:"joinDate" gorm:"not null"`
	LeaveDate    *time.Time `json:"leaveDate"`
	HourlyRate   float64 `json:"hourlyRate"`
	Skills       string  `json:"skills"`
} 