package infra

func (this *TestServicePointHandler) Reset() {
	this.postNumber.Add(this.postNumber.Load() * -1)
	this.putNumber.Add(this.putNumber.Load() * -1)
	this.patchNumber.Add(this.patchNumber.Load() * -1)
	this.getNumber.Add(this.getNumber.Load() * -1)
	this.deleteNumber.Add(this.deleteNumber.Load() * -1)
	this.tr = false
	this.errorMode = false
}

func (this *TestServicePointHandler) Name() string {
	return this.name
}

func (this *TestServicePointHandler) Tr() bool {
	return this.tr
}

func (this *TestServicePointHandler) SetTr(b bool) {
	this.tr = b
}

func (this *TestServicePointHandler) ErrorMode() bool {
	return this.errorMode
}

func (this *TestServicePointHandler) SetErrorMode(b bool) {
	this.errorMode = b
}

func (this *TestServicePointHandler) PostN() int {
	return int(this.postNumber.Load())
}

func (this *TestServicePointHandler) PutN() int {
	return int(this.putNumber.Load())
}

func (this *TestServicePointHandler) PatchN() int {
	return int(this.patchNumber.Load())
}

func (this *TestServicePointHandler) GetN() int {
	return int(this.getNumber.Load())
}

func (this *TestServicePointHandler) DeleteN() int {
	return int(this.deleteNumber.Load())
}

func (this *TestServicePointHandler) FailedN() int {
	return int(this.failedNumber.Load())
}
