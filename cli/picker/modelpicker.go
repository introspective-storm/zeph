package picker

func ModelPicker(height int) Model {
	allowedExt := []string{
		".jls",
		".jld2",
		".bson",
		".pkl",
		".pickle",
		".joblib",
		".pt",
		".pth",
		".h5",
		".hdf5",
		".keras",
	}
	return New(allowedExt, height)
}
