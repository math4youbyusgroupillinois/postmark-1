package postmark

type Bounce struct {
    ID, Type, Tag                        string
    TypeCode                             int
    MessageID, BouncedAT, Details, Email string
    DumpAvailable, Inactive, CanActivate bool
    Subject                              string
}
