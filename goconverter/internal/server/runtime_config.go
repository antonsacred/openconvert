package server

func ConfigureRuntimeLimits(decodedFileSizeBytes int, requestBodyBytes int64, concurrentConversions int) {
	if decodedFileSizeBytes > 0 {
		maxDecodedFileSizeBytes = decodedFileSizeBytes
	}

	if requestBodyBytes > 0 {
		maxRequestBodyBytes = requestBodyBytes
	}

	if concurrentConversions > 0 {
		maxConcurrentConversions = concurrentConversions
	}
}
