// Objective-C API for talking to golang.org/x/mobile/bind/objc/testpkg Go package.
//   gobind -lang=objc golang.org/x/mobile/bind/objc/testpkg
//
// File is generated by gobind. Do not edit.

#include "GoTestpkg.h"
#include <Foundation/Foundation.h>
#include "seq.h"

static NSString* errDomain = @"go.golang.org/x/mobile/bind/objc/testpkg";

@protocol goSeqRefInterface
-(GoSeqRef*) ref;
@end

#define _DESCRIPTOR_ "testpkg"

#define _CALL_BytesAppend_ 1
#define _CALL_CallIError_ 2
#define _CALL_CallIStringError_ 3
#define _CALL_CallSSum_ 4
#define _CALL_CollectS_ 5
#define _CALL_GC_ 6
#define _CALL_Hello_ 7
#define _CALL_Hi_ 8
#define _CALL_Int_ 9
#define _CALL_Multiply_ 10
#define _CALL_NewI_ 11
#define _CALL_NewNode_ 12
#define _CALL_NewS_ 13
#define _CALL_RegisterI_ 14
#define _CALL_ReturnsError_ 15
#define _CALL_Sum_ 16
#define _CALL_UnregisterI_ 17

#define _GO_testpkg_I_DESCRIPTOR_ "go.testpkg.I"
#define _GO_testpkg_I_Error_ (0x10a)
#define _GO_testpkg_I_StringError_ (0x20a)
#define _GO_testpkg_I_Times_ (0x30a)

@interface GoTestpkgI : NSObject <GoTestpkgI> {
}
@property(strong, readonly) id ref;

- (id)initWithRef:(id)ref;
- (BOOL)Error:(BOOL)triggerError error:(NSError**)error;
- (BOOL)StringError:(NSString*)s ret0_:(NSString**)ret0_ error:(NSError**)error;
- (int64_t)Times:(int32_t)v;
@end

@implementation GoTestpkgI {
}

- (id)initWithRef:(id)ref {
	self = [super init];
	if (self) { _ref = ref; }
	return self;
}

- (BOOL)Error:(BOOL)triggerError error:(NSError**)error {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_writeBool(&in_, triggerError);
	go_seq_send(_GO_testpkg_I_DESCRIPTOR_, _GO_testpkg_I_Error_, &in_, &out_);
	NSString* _error = go_seq_readUTF8(&out_);
	if ([_error length] != 0 && error != nil) {
		NSMutableDictionary* details = [NSMutableDictionary dictionary];
		[details setValue:_error forKey:NSLocalizedDescriptionKey];
		*error = [NSError errorWithDomain:errDomain code:1 userInfo:details];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ([_error length] == 0);
}

- (BOOL)StringError:(NSString*)s ret0_:(NSString**)ret0_ error:(NSError**)error {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_writeUTF8(&in_, s);
	go_seq_send(_GO_testpkg_I_DESCRIPTOR_, _GO_testpkg_I_StringError_, &in_, &out_);
	NSString* ret0__val = go_seq_readUTF8(&out_);
	if (ret0_ != NULL) {
		*ret0_ = ret0__val;
	}
	NSString* _error = go_seq_readUTF8(&out_);
	if ([_error length] != 0 && error != nil) {
		NSMutableDictionary* details = [NSMutableDictionary dictionary];
		[details setValue:_error forKey:NSLocalizedDescriptionKey];
		*error = [NSError errorWithDomain:errDomain code:1 userInfo:details];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ([_error length] == 0);
}

- (int64_t)Times:(int32_t)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_writeInt32(&in_, v);
	go_seq_send(_GO_testpkg_I_DESCRIPTOR_, _GO_testpkg_I_Times_, &in_, &out_);
	int64_t ret0_ = go_seq_readInt64(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

@end

static void proxyGoTestpkgI(id obj, int code, GoSeq* in, GoSeq* out) {
	switch (code) {
	case _GO_testpkg_I_Error_: {
		id<GoTestpkgI> o = (id<GoTestpkgI>)(obj);
		BOOL triggerError = go_seq_readBool(in);
		NSError* error = NULL;
		BOOL returnVal = [o Error:triggerError error:&error];
		if (returnVal) {
			go_seq_writeUTF8(out, NULL);
		} else {
			NSString* errorDesc = [error localizedDescription];
			if (errorDesc == NULL || errorDesc.length == 0) {
				errorDesc = @"gobind: unknown error";
			}
			go_seq_writeUTF8(out, errorDesc);
		}
	} break;
	case _GO_testpkg_I_StringError_: {
		id<GoTestpkgI> o = (id<GoTestpkgI>)(obj);
		NSString* s = go_seq_readUTF8(in);
		NSString* ret0_;
		NSError* error = NULL;
		BOOL returnVal = [o StringError:s ret0_:&ret0_ error:&error];
		go_seq_writeUTF8(out, ret0_);
		if (returnVal) {
			go_seq_writeUTF8(out, NULL);
		} else {
			NSString* errorDesc = [error localizedDescription];
			if (errorDesc == NULL || errorDesc.length == 0) {
				errorDesc = @"gobind: unknown error";
			}
			go_seq_writeUTF8(out, errorDesc);
		}
	} break;
	case _GO_testpkg_I_Times_: {
		id<GoTestpkgI> o = (id<GoTestpkgI>)(obj);
		int32_t v = go_seq_readInt32(in);
		int64_t returnVal = [o Times:v];
		go_seq_writeInt64(out, returnVal);
	} break;
	default:
		NSLog(@"unknown code %x for _GO_testpkg_I_DESCRIPTOR_", code);
	}
}

#define _GO_testpkg_Node_DESCRIPTOR_ "go.testpkg.Node"
#define _GO_testpkg_Node_FIELD_V_GET_ (0x00f)
#define _GO_testpkg_Node_FIELD_V_SET_ (0x01f)
#define _GO_testpkg_Node_FIELD_Err_GET_ (0x10f)
#define _GO_testpkg_Node_FIELD_Err_SET_ (0x11f)

@implementation GoTestpkgNode {
}

- (id)initWithRef:(id)ref {
	self = [super init];
	if (self) { _ref = ref; }
	return self;
}

- (NSString*)V {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_send(_GO_testpkg_Node_DESCRIPTOR_, _GO_testpkg_Node_FIELD_V_GET_, &in_, &out_);
	NSString* ret_ = go_seq_readUTF8(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret_;
}

- (void)setV:(NSString*)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_writeUTF8(&in_, v);
	go_seq_send(_GO_testpkg_Node_DESCRIPTOR_, _GO_testpkg_Node_FIELD_V_SET_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

- (NSString*)Err {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_send(_GO_testpkg_Node_DESCRIPTOR_, _GO_testpkg_Node_FIELD_Err_GET_, &in_, &out_);
	NSString* ret_ = go_seq_readUTF8(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret_;
}

- (void)setErr:(NSString*)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_writeUTF8(&in_, v);
	go_seq_send(_GO_testpkg_Node_DESCRIPTOR_, _GO_testpkg_Node_FIELD_Err_SET_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

@end

#define _GO_testpkg_S_DESCRIPTOR_ "go.testpkg.S"
#define _GO_testpkg_S_FIELD_X_GET_ (0x00f)
#define _GO_testpkg_S_FIELD_X_SET_ (0x01f)
#define _GO_testpkg_S_FIELD_Y_GET_ (0x10f)
#define _GO_testpkg_S_FIELD_Y_SET_ (0x11f)
#define _GO_testpkg_S_Sum_ (0x00c)
#define _GO_testpkg_S_TryTwoStrings_ (0x10c)

@implementation GoTestpkgS {
}

- (id)initWithRef:(id)ref {
	self = [super init];
	if (self) { _ref = ref; }
	return self;
}

- (double)X {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_send(_GO_testpkg_S_DESCRIPTOR_, _GO_testpkg_S_FIELD_X_GET_, &in_, &out_);
	double ret_ = go_seq_readFloat64(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret_;
}

- (void)setX:(double)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_writeFloat64(&in_, v);
	go_seq_send(_GO_testpkg_S_DESCRIPTOR_, _GO_testpkg_S_FIELD_X_SET_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

- (double)Y {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_send(_GO_testpkg_S_DESCRIPTOR_, _GO_testpkg_S_FIELD_Y_GET_, &in_, &out_);
	double ret_ = go_seq_readFloat64(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret_;
}

- (void)setY:(double)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_writeFloat64(&in_, v);
	go_seq_send(_GO_testpkg_S_DESCRIPTOR_, _GO_testpkg_S_FIELD_Y_SET_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

- (double)Sum {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_send(_GO_testpkg_S_DESCRIPTOR_, _GO_testpkg_S_Sum_, &in_, &out_);
	double ret0_ = go_seq_readFloat64(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

- (NSString*)TryTwoStrings:(NSString*)first second:(NSString*)second {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeRef(&in_, self.ref);
	go_seq_writeUTF8(&in_, first);
	go_seq_writeUTF8(&in_, second);
	go_seq_send(_GO_testpkg_S_DESCRIPTOR_, _GO_testpkg_S_TryTwoStrings_, &in_, &out_);
	NSString* ret0_ = go_seq_readUTF8(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

@end

const BOOL GoTestpkgABool = YES;
const double GoTestpkgAFloat = 0.12345;
NSString* const GoTestpkgAString = @"a string";
const int64_t GoTestpkgAnInt = 7LL;
const double GoTestpkgLog2E = 1.4426950408889634;
const float GoTestpkgMaxFloat32 = 3.4028234663852886e+38;
const double GoTestpkgMaxFloat64 = 1.7976931348623157e+308;
const int32_t GoTestpkgMaxInt32 = 2147483647;
const int64_t GoTestpkgMaxInt64 = 9223372036854775807LL;
const int32_t GoTestpkgMinInt32 = -2147483648;
const int64_t GoTestpkgMinInt64 = -9223372036854775807LL-1;
const float GoTestpkgSmallestNonzeroFloat32 = 0;
const double GoTestpkgSmallestNonzeroFloat64 = 5e-324;

@implementation GoTestpkg
+ (void) setIntVar:(int)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeInt(&in_, v);
	go_seq_send("testpkg.IntVar", 1, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

+ (int) IntVar {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_send("testpkg.IntVar", 2, &in_, &out_);
	int ret = go_seq_readInt(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret;
}

+ (void) setInterfaceVar:(id<GoTestpkgI>)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	if ([(id<NSObject>)(v) isKindOfClass:[GoTestpkgI class]]) {
		id<goSeqRefInterface> v_proxy = (id<goSeqRefInterface>)(v);
		go_seq_writeRef(&in_, v_proxy.ref);
	} else {
		go_seq_writeObjcRef(&in_, v);
	}
	go_seq_send("testpkg.InterfaceVar", 1, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

+ (id<GoTestpkgI>) InterfaceVar {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_send("testpkg.InterfaceVar", 2, &in_, &out_);
	GoSeqRef* ret_ref = go_seq_readRef(&out_);
	id<GoTestpkgI> ret = ret_ref.obj;
	if (ret == NULL) {
		ret = [[GoTestpkgI alloc] initWithRef:ret_ref];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret;
}

+ (void) setStringVar:(NSString*)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeUTF8(&in_, v);
	go_seq_send("testpkg.StringVar", 1, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

+ (NSString*) StringVar {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_send("testpkg.StringVar", 2, &in_, &out_);
	NSString* ret = go_seq_readUTF8(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret;
}

+ (void) setStructVar:(GoTestpkgNode*)v {
	GoSeq in_ = {};
	GoSeq out_ = {};
	if ([(id<NSObject>)(v) isKindOfClass:[GoTestpkgNode class]]) {
		id<goSeqRefInterface> v_proxy = (id<goSeqRefInterface>)(v);
		go_seq_writeRef(&in_, v_proxy.ref);
	} else {
		go_seq_writeObjcRef(&in_, v);
	}
	go_seq_send("testpkg.StructVar", 1, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

+ (GoTestpkgNode*) StructVar {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_send("testpkg.StructVar", 2, &in_, &out_);
	GoSeqRef* ret_ref = go_seq_readRef(&out_);
	GoTestpkgNode* ret = ret_ref.obj;
	if (ret == NULL) {
		ret = [[GoTestpkgNode alloc] initWithRef:ret_ref];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret;
}

@end

NSData* GoTestpkgBytesAppend(NSData* a, NSData* b) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeByteArray(&in_, a);
	go_seq_writeByteArray(&in_, b);
	go_seq_send(_DESCRIPTOR_, _CALL_BytesAppend_, &in_, &out_);
	NSData* ret0_ = go_seq_readByteArray(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

BOOL GoTestpkgCallIError(id<GoTestpkgI> i, BOOL triggerError, NSError** error) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	if ([(id<NSObject>)(i) isKindOfClass:[GoTestpkgI class]]) {
		id<goSeqRefInterface> i_proxy = (id<goSeqRefInterface>)(i);
		go_seq_writeRef(&in_, i_proxy.ref);
	} else {
		go_seq_writeObjcRef(&in_, i);
	}
	go_seq_writeBool(&in_, triggerError);
	go_seq_send(_DESCRIPTOR_, _CALL_CallIError_, &in_, &out_);
	NSString* _error = go_seq_readUTF8(&out_);
	if ([_error length] != 0 && error != nil) {
		NSMutableDictionary* details = [NSMutableDictionary dictionary];
		[details setValue:_error forKey:NSLocalizedDescriptionKey];
		*error = [NSError errorWithDomain:errDomain code:1 userInfo:details];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ([_error length] == 0);
}

BOOL GoTestpkgCallIStringError(id<GoTestpkgI> i, NSString* s, NSString** ret0_, NSError** error) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	if ([(id<NSObject>)(i) isKindOfClass:[GoTestpkgI class]]) {
		id<goSeqRefInterface> i_proxy = (id<goSeqRefInterface>)(i);
		go_seq_writeRef(&in_, i_proxy.ref);
	} else {
		go_seq_writeObjcRef(&in_, i);
	}
	go_seq_writeUTF8(&in_, s);
	go_seq_send(_DESCRIPTOR_, _CALL_CallIStringError_, &in_, &out_);
	NSString* ret0__val = go_seq_readUTF8(&out_);
	if (ret0_ != NULL) {
		*ret0_ = ret0__val;
	}
	NSString* _error = go_seq_readUTF8(&out_);
	if ([_error length] != 0 && error != nil) {
		NSMutableDictionary* details = [NSMutableDictionary dictionary];
		[details setValue:_error forKey:NSLocalizedDescriptionKey];
		*error = [NSError errorWithDomain:errDomain code:1 userInfo:details];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ([_error length] == 0);
}

double GoTestpkgCallSSum(GoTestpkgS* s) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	if ([(id<NSObject>)(s) isKindOfClass:[GoTestpkgS class]]) {
		id<goSeqRefInterface> s_proxy = (id<goSeqRefInterface>)(s);
		go_seq_writeRef(&in_, s_proxy.ref);
	} else {
		go_seq_writeObjcRef(&in_, s);
	}
	go_seq_send(_DESCRIPTOR_, _CALL_CallSSum_, &in_, &out_);
	double ret0_ = go_seq_readFloat64(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

int GoTestpkgCollectS(int want, int timeoutSec) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeInt(&in_, want);
	go_seq_writeInt(&in_, timeoutSec);
	go_seq_send(_DESCRIPTOR_, _CALL_CollectS_, &in_, &out_);
	int ret0_ = go_seq_readInt(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

void GoTestpkgGC() {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_send(_DESCRIPTOR_, _CALL_GC_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

NSString* GoTestpkgHello(NSString* s) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeUTF8(&in_, s);
	go_seq_send(_DESCRIPTOR_, _CALL_Hello_, &in_, &out_);
	NSString* ret0_ = go_seq_readUTF8(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

void GoTestpkgHi() {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_send(_DESCRIPTOR_, _CALL_Hi_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

void GoTestpkgInt(int32_t x) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeInt32(&in_, x);
	go_seq_send(_DESCRIPTOR_, _CALL_Int_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

int64_t GoTestpkgMultiply(int32_t idx, int32_t val) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeInt32(&in_, idx);
	go_seq_writeInt32(&in_, val);
	go_seq_send(_DESCRIPTOR_, _CALL_Multiply_, &in_, &out_);
	int64_t ret0_ = go_seq_readInt64(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

id<GoTestpkgI> GoTestpkgNewI() {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_send(_DESCRIPTOR_, _CALL_NewI_, &in_, &out_);
	GoSeqRef* ret0__ref = go_seq_readRef(&out_);
	id<GoTestpkgI> ret0_ = ret0__ref.obj;
	if (ret0_ == NULL) {
		ret0_ = [[GoTestpkgI alloc] initWithRef:ret0__ref];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

GoTestpkgNode* GoTestpkgNewNode(NSString* name) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeUTF8(&in_, name);
	go_seq_send(_DESCRIPTOR_, _CALL_NewNode_, &in_, &out_);
	GoSeqRef* ret0__ref = go_seq_readRef(&out_);
	GoTestpkgNode* ret0_ = ret0__ref.obj;
	if (ret0_ == NULL) {
		ret0_ = [[GoTestpkgNode alloc] initWithRef:ret0__ref];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

GoTestpkgS* GoTestpkgNewS(double x, double y) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeFloat64(&in_, x);
	go_seq_writeFloat64(&in_, y);
	go_seq_send(_DESCRIPTOR_, _CALL_NewS_, &in_, &out_);
	GoSeqRef* ret0__ref = go_seq_readRef(&out_);
	GoTestpkgS* ret0_ = ret0__ref.obj;
	if (ret0_ == NULL) {
		ret0_ = [[GoTestpkgS alloc] initWithRef:ret0__ref];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

void GoTestpkgRegisterI(int32_t idx, id<GoTestpkgI> i) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeInt32(&in_, idx);
	if ([(id<NSObject>)(i) isKindOfClass:[GoTestpkgI class]]) {
		id<goSeqRefInterface> i_proxy = (id<goSeqRefInterface>)(i);
		go_seq_writeRef(&in_, i_proxy.ref);
	} else {
		go_seq_writeObjcRef(&in_, i);
	}
	go_seq_send(_DESCRIPTOR_, _CALL_RegisterI_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

BOOL GoTestpkgReturnsError(BOOL b, NSString** ret0_, NSError** error) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeBool(&in_, b);
	go_seq_send(_DESCRIPTOR_, _CALL_ReturnsError_, &in_, &out_);
	NSString* ret0__val = go_seq_readUTF8(&out_);
	if (ret0_ != NULL) {
		*ret0_ = ret0__val;
	}
	NSString* _error = go_seq_readUTF8(&out_);
	if ([_error length] != 0 && error != nil) {
		NSMutableDictionary* details = [NSMutableDictionary dictionary];
		[details setValue:_error forKey:NSLocalizedDescriptionKey];
		*error = [NSError errorWithDomain:errDomain code:1 userInfo:details];
	}
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ([_error length] == 0);
}

int64_t GoTestpkgSum(int64_t x, int64_t y) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeInt64(&in_, x);
	go_seq_writeInt64(&in_, y);
	go_seq_send(_DESCRIPTOR_, _CALL_Sum_, &in_, &out_);
	int64_t ret0_ = go_seq_readInt64(&out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
	return ret0_;
}

void GoTestpkgUnregisterI(int32_t idx) {
	GoSeq in_ = {};
	GoSeq out_ = {};
	go_seq_writeInt32(&in_, idx);
	go_seq_send(_DESCRIPTOR_, _CALL_UnregisterI_, &in_, &out_);
	go_seq_free(&in_);
	go_seq_free(&out_);
}

__attribute__((constructor)) static void init() {
	go_seq_register_proxy("go.testpkg.I", proxyGoTestpkgI);
}
