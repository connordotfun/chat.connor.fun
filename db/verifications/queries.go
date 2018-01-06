package verifications


const (
	createIfNotExistsQuery = `
		CREATE IF NOT EXISTS verification_codes (
			code VARCHAR(50) PRIMARY KEY,
			purpose VARCHAR(50) NOT NULL,
			user_id UUID NOT NULL REFERENCES users (id),
			valid Boolean NOT NULL,
			exp_date Integer
		);
	`

	insertCodeQuery = `
		INSERT INTO verification_codes (code, purpose, user_id, valid, exp_date) VALUES (
			code = :code, purpose = :purpose, user_id = :user_id, valid = :valid, exp_date = :exp_date
		);
	`

	invalidateCodeQuery = `
		UPDATE verification_codes
			SET valid = false
		WHERE code = :code;

	`

	selectByCodeQuery = `
		SELECT * FROM verification_codes
			WHERE code = :code;
	`
)
