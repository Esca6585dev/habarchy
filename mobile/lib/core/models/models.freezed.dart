// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint, type=warning, deprecated_member_use, deprecated_member_use_from_same_package
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'models.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$ApiError {

 String get code; String get message; Map<String, dynamic>? get details;
/// Create a copy of ApiError
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ApiErrorCopyWith<ApiError> get copyWith => _$ApiErrorCopyWithImpl<ApiError>(this as ApiError, _$identity);

  /// Serializes this ApiError to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as ApiError;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ApiError&&(identical(other.code, _this.code) || other.code == _this.code)&&(identical(other.message, _this.message) || other.message == _this.message)&&const DeepCollectionEquality().equals(other.details, _this.details));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as ApiError;
  return Object.hash(runtimeType,_this.code,_this.message,const DeepCollectionEquality().hash(_this.details));
}

@override
String toString() {
  final _this = this as ApiError;
  return 'ApiError(code: ${_this.code}, message: ${_this.message}, details: ${_this.details})';
}


}

/// @nodoc
abstract mixin class $ApiErrorCopyWith<$Res>  {
  factory $ApiErrorCopyWith(ApiError value, $Res Function(ApiError) _then) = _$ApiErrorCopyWithImpl;
@useResult
$Res call({
 String code, String message, Map<String, dynamic>? details
});




}
/// @nodoc
class _$ApiErrorCopyWithImpl<$Res>
    implements $ApiErrorCopyWith<$Res> {
  _$ApiErrorCopyWithImpl(this._self, this._then);

  final ApiError _self;
  final $Res Function(ApiError) _then;

/// Create a copy of ApiError
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? code = null,Object? message = null,Object? details = freezed,}) {
  return _then(ApiError(
code: null == code ? _self.code : code // ignore: cast_nullable_to_non_nullable
as String,message: null == message ? _self.message : message // ignore: cast_nullable_to_non_nullable
as String,details: freezed == details ? _self.details : details // ignore: cast_nullable_to_non_nullable
as Map<String, dynamic>?,
  ));
}

}


/// Adds pattern-matching-related methods to [ApiError].
extension ApiErrorPatterns on ApiError {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ApiError value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ApiError() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ApiError value)  $default,){
final _that = this;
switch (_that) {
case _ApiError():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ApiError value)?  $default,){
final _that = this;
switch (_that) {
case _ApiError() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String code,  String message,  Map<String, dynamic>? details)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ApiError() when $default != null:
return $default(_that.code,_that.message,_that.details);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String code,  String message,  Map<String, dynamic>? details)  $default,) {final _that = this;
switch (_that) {
case _ApiError():
return $default(_that.code,_that.message,_that.details);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String code,  String message,  Map<String, dynamic>? details)?  $default,) {final _that = this;
switch (_that) {
case _ApiError() when $default != null:
return $default(_that.code,_that.message,_that.details);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ApiError implements ApiError {
  const _ApiError({required this.code, required this.message,  Map<String, dynamic>? details}): _details = details;
  factory _ApiError.fromJson(Map<String, dynamic> json) => _$ApiErrorFromJson(json);

@override final  String code;
@override final  String message;
 final  Map<String, dynamic>? _details;
@override Map<String, dynamic>? get details {
  final value = _details;
  if (value == null) return null;
  if (_details is EqualUnmodifiableMapView) return _details;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableMapView(value);
}


/// Create a copy of ApiError
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ApiErrorCopyWith<_ApiError> get copyWith => __$ApiErrorCopyWithImpl<_ApiError>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ApiErrorToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _ApiError&&(identical(other.code, code) || other.code == code)&&(identical(other.message, message) || other.message == message)&&const DeepCollectionEquality().equals(other.details, _details));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,code,message,const DeepCollectionEquality().hash(_details));
}

@override
String toString() {
    return 'ApiError(code: $code, message: $message, details: $details)';
}


}

/// @nodoc
abstract mixin class _$ApiErrorCopyWith<$Res> implements $ApiErrorCopyWith<$Res> {
  factory _$ApiErrorCopyWith(_ApiError value, $Res Function(_ApiError) _then) = __$ApiErrorCopyWithImpl;
@override @useResult
$Res call({
 String code, String message, Map<String, dynamic>? details
});




}
/// @nodoc
class __$ApiErrorCopyWithImpl<$Res>
    implements _$ApiErrorCopyWith<$Res> {
  __$ApiErrorCopyWithImpl(this._self, this._then);

  final _ApiError _self;
  final $Res Function(_ApiError) _then;

/// Create a copy of ApiError
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? code = null,Object? message = null,Object? details = freezed,}) {
  return _then(_ApiError(
code: null == code ? _self.code : code // ignore: cast_nullable_to_non_nullable
as String,message: null == message ? _self.message : message // ignore: cast_nullable_to_non_nullable
as String,details: freezed == details ? _self._details : details // ignore: cast_nullable_to_non_nullable
as Map<String, dynamic>?,
  ));
}


}


/// @nodoc
mixin _$Tokens {

@JsonKey(name: 'access_token') String get accessToken;@JsonKey(name: 'refresh_token') String get refreshToken;@JsonKey(name: 'expires_in') int get expiresIn;
/// Create a copy of Tokens
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TokensCopyWith<Tokens> get copyWith => _$TokensCopyWithImpl<Tokens>(this as Tokens, _$identity);

  /// Serializes this Tokens to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Tokens;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Tokens&&(identical(other.accessToken, _this.accessToken) || other.accessToken == _this.accessToken)&&(identical(other.refreshToken, _this.refreshToken) || other.refreshToken == _this.refreshToken)&&(identical(other.expiresIn, _this.expiresIn) || other.expiresIn == _this.expiresIn));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Tokens;
  return Object.hash(runtimeType,_this.accessToken,_this.refreshToken,_this.expiresIn);
}

@override
String toString() {
  final _this = this as Tokens;
  return 'Tokens(accessToken: ${_this.accessToken}, refreshToken: ${_this.refreshToken}, expiresIn: ${_this.expiresIn})';
}


}

/// @nodoc
abstract mixin class $TokensCopyWith<$Res>  {
  factory $TokensCopyWith(Tokens value, $Res Function(Tokens) _then) = _$TokensCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'access_token') String accessToken,@JsonKey(name: 'refresh_token') String refreshToken,@JsonKey(name: 'expires_in') int expiresIn
});




}
/// @nodoc
class _$TokensCopyWithImpl<$Res>
    implements $TokensCopyWith<$Res> {
  _$TokensCopyWithImpl(this._self, this._then);

  final Tokens _self;
  final $Res Function(Tokens) _then;

/// Create a copy of Tokens
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? accessToken = null,Object? refreshToken = null,Object? expiresIn = null,}) {
  return _then(Tokens(
accessToken: null == accessToken ? _self.accessToken : accessToken // ignore: cast_nullable_to_non_nullable
as String,refreshToken: null == refreshToken ? _self.refreshToken : refreshToken // ignore: cast_nullable_to_non_nullable
as String,expiresIn: null == expiresIn ? _self.expiresIn : expiresIn // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [Tokens].
extension TokensPatterns on Tokens {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Tokens value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Tokens() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Tokens value)  $default,){
final _that = this;
switch (_that) {
case _Tokens():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Tokens value)?  $default,){
final _that = this;
switch (_that) {
case _Tokens() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'access_token')  String accessToken, @JsonKey(name: 'refresh_token')  String refreshToken, @JsonKey(name: 'expires_in')  int expiresIn)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Tokens() when $default != null:
return $default(_that.accessToken,_that.refreshToken,_that.expiresIn);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'access_token')  String accessToken, @JsonKey(name: 'refresh_token')  String refreshToken, @JsonKey(name: 'expires_in')  int expiresIn)  $default,) {final _that = this;
switch (_that) {
case _Tokens():
return $default(_that.accessToken,_that.refreshToken,_that.expiresIn);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'access_token')  String accessToken, @JsonKey(name: 'refresh_token')  String refreshToken, @JsonKey(name: 'expires_in')  int expiresIn)?  $default,) {final _that = this;
switch (_that) {
case _Tokens() when $default != null:
return $default(_that.accessToken,_that.refreshToken,_that.expiresIn);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Tokens implements Tokens {
  const _Tokens({@JsonKey(name: 'access_token') required this.accessToken, @JsonKey(name: 'refresh_token') required this.refreshToken, @JsonKey(name: 'expires_in') this.expiresIn = 900});
  factory _Tokens.fromJson(Map<String, dynamic> json) => _$TokensFromJson(json);

@override@JsonKey(name: 'access_token') final  String accessToken;
@override@JsonKey(name: 'refresh_token') final  String refreshToken;
@override@JsonKey(name: 'expires_in') final  int expiresIn;

/// Create a copy of Tokens
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TokensCopyWith<_Tokens> get copyWith => __$TokensCopyWithImpl<_Tokens>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$TokensToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Tokens&&(identical(other.accessToken, accessToken) || other.accessToken == accessToken)&&(identical(other.refreshToken, refreshToken) || other.refreshToken == refreshToken)&&(identical(other.expiresIn, expiresIn) || other.expiresIn == expiresIn));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,accessToken,refreshToken,expiresIn);
}

@override
String toString() {
    return 'Tokens(accessToken: $accessToken, refreshToken: $refreshToken, expiresIn: $expiresIn)';
}


}

/// @nodoc
abstract mixin class _$TokensCopyWith<$Res> implements $TokensCopyWith<$Res> {
  factory _$TokensCopyWith(_Tokens value, $Res Function(_Tokens) _then) = __$TokensCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'access_token') String accessToken,@JsonKey(name: 'refresh_token') String refreshToken,@JsonKey(name: 'expires_in') int expiresIn
});




}
/// @nodoc
class __$TokensCopyWithImpl<$Res>
    implements _$TokensCopyWith<$Res> {
  __$TokensCopyWithImpl(this._self, this._then);

  final _Tokens _self;
  final $Res Function(_Tokens) _then;

/// Create a copy of Tokens
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? accessToken = null,Object? refreshToken = null,Object? expiresIn = null,}) {
  return _then(_Tokens(
accessToken: null == accessToken ? _self.accessToken : accessToken // ignore: cast_nullable_to_non_nullable
as String,refreshToken: null == refreshToken ? _self.refreshToken : refreshToken // ignore: cast_nullable_to_non_nullable
as String,expiresIn: null == expiresIn ? _self.expiresIn : expiresIn // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$User {

 String get id; String get email;@JsonKey(name: 'full_name') String get fullName;@JsonKey(name: 'totp_enabled') bool get totpEnabled;
/// Create a copy of User
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$UserCopyWith<User> get copyWith => _$UserCopyWithImpl<User>(this as User, _$identity);

  /// Serializes this User to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as User;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is User&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.email, _this.email) || other.email == _this.email)&&(identical(other.fullName, _this.fullName) || other.fullName == _this.fullName)&&(identical(other.totpEnabled, _this.totpEnabled) || other.totpEnabled == _this.totpEnabled));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as User;
  return Object.hash(runtimeType,_this.id,_this.email,_this.fullName,_this.totpEnabled);
}

@override
String toString() {
  final _this = this as User;
  return 'User(id: ${_this.id}, email: ${_this.email}, fullName: ${_this.fullName}, totpEnabled: ${_this.totpEnabled})';
}


}

/// @nodoc
abstract mixin class $UserCopyWith<$Res>  {
  factory $UserCopyWith(User value, $Res Function(User) _then) = _$UserCopyWithImpl;
@useResult
$Res call({
 String id, String email,@JsonKey(name: 'full_name') String fullName,@JsonKey(name: 'totp_enabled') bool totpEnabled
});




}
/// @nodoc
class _$UserCopyWithImpl<$Res>
    implements $UserCopyWith<$Res> {
  _$UserCopyWithImpl(this._self, this._then);

  final User _self;
  final $Res Function(User) _then;

/// Create a copy of User
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? email = null,Object? fullName = null,Object? totpEnabled = null,}) {
  return _then(User(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,email: null == email ? _self.email : email // ignore: cast_nullable_to_non_nullable
as String,fullName: null == fullName ? _self.fullName : fullName // ignore: cast_nullable_to_non_nullable
as String,totpEnabled: null == totpEnabled ? _self.totpEnabled : totpEnabled // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

}


/// Adds pattern-matching-related methods to [User].
extension UserPatterns on User {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _User value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _User() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _User value)  $default,){
final _that = this;
switch (_that) {
case _User():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _User value)?  $default,){
final _that = this;
switch (_that) {
case _User() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String email, @JsonKey(name: 'full_name')  String fullName, @JsonKey(name: 'totp_enabled')  bool totpEnabled)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _User() when $default != null:
return $default(_that.id,_that.email,_that.fullName,_that.totpEnabled);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String email, @JsonKey(name: 'full_name')  String fullName, @JsonKey(name: 'totp_enabled')  bool totpEnabled)  $default,) {final _that = this;
switch (_that) {
case _User():
return $default(_that.id,_that.email,_that.fullName,_that.totpEnabled);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String email, @JsonKey(name: 'full_name')  String fullName, @JsonKey(name: 'totp_enabled')  bool totpEnabled)?  $default,) {final _that = this;
switch (_that) {
case _User() when $default != null:
return $default(_that.id,_that.email,_that.fullName,_that.totpEnabled);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _User implements User {
  const _User({required this.id, required this.email, @JsonKey(name: 'full_name') this.fullName = '', @JsonKey(name: 'totp_enabled') this.totpEnabled = false});
  factory _User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);

@override final  String id;
@override final  String email;
@override@JsonKey(name: 'full_name') final  String fullName;
@override@JsonKey(name: 'totp_enabled') final  bool totpEnabled;

/// Create a copy of User
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$UserCopyWith<_User> get copyWith => __$UserCopyWithImpl<_User>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$UserToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _User&&(identical(other.id, id) || other.id == id)&&(identical(other.email, email) || other.email == email)&&(identical(other.fullName, fullName) || other.fullName == fullName)&&(identical(other.totpEnabled, totpEnabled) || other.totpEnabled == totpEnabled));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,email,fullName,totpEnabled);
}

@override
String toString() {
    return 'User(id: $id, email: $email, fullName: $fullName, totpEnabled: $totpEnabled)';
}


}

/// @nodoc
abstract mixin class _$UserCopyWith<$Res> implements $UserCopyWith<$Res> {
  factory _$UserCopyWith(_User value, $Res Function(_User) _then) = __$UserCopyWithImpl;
@override @useResult
$Res call({
 String id, String email,@JsonKey(name: 'full_name') String fullName,@JsonKey(name: 'totp_enabled') bool totpEnabled
});




}
/// @nodoc
class __$UserCopyWithImpl<$Res>
    implements _$UserCopyWith<$Res> {
  __$UserCopyWithImpl(this._self, this._then);

  final _User _self;
  final $Res Function(_User) _then;

/// Create a copy of User
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? email = null,Object? fullName = null,Object? totpEnabled = null,}) {
  return _then(_User(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,email: null == email ? _self.email : email // ignore: cast_nullable_to_non_nullable
as String,fullName: null == fullName ? _self.fullName : fullName // ignore: cast_nullable_to_non_nullable
as String,totpEnabled: null == totpEnabled ? _self.totpEnabled : totpEnabled // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}


}


/// @nodoc
mixin _$Project {

 String get id; String get name; String get slug; String get status; String get role;@JsonKey(name: 'daily_quota') int get dailyQuota;@JsonKey(name: 'monthly_quota') int get monthlyQuota;@JsonKey(name: 'default_locale') String get defaultLocale;
/// Create a copy of Project
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ProjectCopyWith<Project> get copyWith => _$ProjectCopyWithImpl<Project>(this as Project, _$identity);

  /// Serializes this Project to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Project;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Project&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.name, _this.name) || other.name == _this.name)&&(identical(other.slug, _this.slug) || other.slug == _this.slug)&&(identical(other.status, _this.status) || other.status == _this.status)&&(identical(other.role, _this.role) || other.role == _this.role)&&(identical(other.dailyQuota, _this.dailyQuota) || other.dailyQuota == _this.dailyQuota)&&(identical(other.monthlyQuota, _this.monthlyQuota) || other.monthlyQuota == _this.monthlyQuota)&&(identical(other.defaultLocale, _this.defaultLocale) || other.defaultLocale == _this.defaultLocale));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Project;
  return Object.hash(runtimeType,_this.id,_this.name,_this.slug,_this.status,_this.role,_this.dailyQuota,_this.monthlyQuota,_this.defaultLocale);
}

@override
String toString() {
  final _this = this as Project;
  return 'Project(id: ${_this.id}, name: ${_this.name}, slug: ${_this.slug}, status: ${_this.status}, role: ${_this.role}, dailyQuota: ${_this.dailyQuota}, monthlyQuota: ${_this.monthlyQuota}, defaultLocale: ${_this.defaultLocale})';
}


}

/// @nodoc
abstract mixin class $ProjectCopyWith<$Res>  {
  factory $ProjectCopyWith(Project value, $Res Function(Project) _then) = _$ProjectCopyWithImpl;
@useResult
$Res call({
 String id, String name, String slug, String status, String role,@JsonKey(name: 'daily_quota') int dailyQuota,@JsonKey(name: 'monthly_quota') int monthlyQuota,@JsonKey(name: 'default_locale') String defaultLocale
});




}
/// @nodoc
class _$ProjectCopyWithImpl<$Res>
    implements $ProjectCopyWith<$Res> {
  _$ProjectCopyWithImpl(this._self, this._then);

  final Project _self;
  final $Res Function(Project) _then;

/// Create a copy of Project
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? name = null,Object? slug = null,Object? status = null,Object? role = null,Object? dailyQuota = null,Object? monthlyQuota = null,Object? defaultLocale = null,}) {
  return _then(Project(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,slug: null == slug ? _self.slug : slug // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,role: null == role ? _self.role : role // ignore: cast_nullable_to_non_nullable
as String,dailyQuota: null == dailyQuota ? _self.dailyQuota : dailyQuota // ignore: cast_nullable_to_non_nullable
as int,monthlyQuota: null == monthlyQuota ? _self.monthlyQuota : monthlyQuota // ignore: cast_nullable_to_non_nullable
as int,defaultLocale: null == defaultLocale ? _self.defaultLocale : defaultLocale // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [Project].
extension ProjectPatterns on Project {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Project value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Project() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Project value)  $default,){
final _that = this;
switch (_that) {
case _Project():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Project value)?  $default,){
final _that = this;
switch (_that) {
case _Project() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String name,  String slug,  String status,  String role, @JsonKey(name: 'daily_quota')  int dailyQuota, @JsonKey(name: 'monthly_quota')  int monthlyQuota, @JsonKey(name: 'default_locale')  String defaultLocale)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Project() when $default != null:
return $default(_that.id,_that.name,_that.slug,_that.status,_that.role,_that.dailyQuota,_that.monthlyQuota,_that.defaultLocale);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String name,  String slug,  String status,  String role, @JsonKey(name: 'daily_quota')  int dailyQuota, @JsonKey(name: 'monthly_quota')  int monthlyQuota, @JsonKey(name: 'default_locale')  String defaultLocale)  $default,) {final _that = this;
switch (_that) {
case _Project():
return $default(_that.id,_that.name,_that.slug,_that.status,_that.role,_that.dailyQuota,_that.monthlyQuota,_that.defaultLocale);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String name,  String slug,  String status,  String role, @JsonKey(name: 'daily_quota')  int dailyQuota, @JsonKey(name: 'monthly_quota')  int monthlyQuota, @JsonKey(name: 'default_locale')  String defaultLocale)?  $default,) {final _that = this;
switch (_that) {
case _Project() when $default != null:
return $default(_that.id,_that.name,_that.slug,_that.status,_that.role,_that.dailyQuota,_that.monthlyQuota,_that.defaultLocale);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Project implements Project {
  const _Project({required this.id, required this.name, required this.slug, this.status = 'active', this.role = 'viewer', @JsonKey(name: 'daily_quota') this.dailyQuota = 0, @JsonKey(name: 'monthly_quota') this.monthlyQuota = 0, @JsonKey(name: 'default_locale') this.defaultLocale = 'tk'});
  factory _Project.fromJson(Map<String, dynamic> json) => _$ProjectFromJson(json);

@override final  String id;
@override final  String name;
@override final  String slug;
@override@JsonKey() final  String status;
@override@JsonKey() final  String role;
@override@JsonKey(name: 'daily_quota') final  int dailyQuota;
@override@JsonKey(name: 'monthly_quota') final  int monthlyQuota;
@override@JsonKey(name: 'default_locale') final  String defaultLocale;

/// Create a copy of Project
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ProjectCopyWith<_Project> get copyWith => __$ProjectCopyWithImpl<_Project>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ProjectToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Project&&(identical(other.id, id) || other.id == id)&&(identical(other.name, name) || other.name == name)&&(identical(other.slug, slug) || other.slug == slug)&&(identical(other.status, status) || other.status == status)&&(identical(other.role, role) || other.role == role)&&(identical(other.dailyQuota, dailyQuota) || other.dailyQuota == dailyQuota)&&(identical(other.monthlyQuota, monthlyQuota) || other.monthlyQuota == monthlyQuota)&&(identical(other.defaultLocale, defaultLocale) || other.defaultLocale == defaultLocale));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,name,slug,status,role,dailyQuota,monthlyQuota,defaultLocale);
}

@override
String toString() {
    return 'Project(id: $id, name: $name, slug: $slug, status: $status, role: $role, dailyQuota: $dailyQuota, monthlyQuota: $monthlyQuota, defaultLocale: $defaultLocale)';
}


}

/// @nodoc
abstract mixin class _$ProjectCopyWith<$Res> implements $ProjectCopyWith<$Res> {
  factory _$ProjectCopyWith(_Project value, $Res Function(_Project) _then) = __$ProjectCopyWithImpl;
@override @useResult
$Res call({
 String id, String name, String slug, String status, String role,@JsonKey(name: 'daily_quota') int dailyQuota,@JsonKey(name: 'monthly_quota') int monthlyQuota,@JsonKey(name: 'default_locale') String defaultLocale
});




}
/// @nodoc
class __$ProjectCopyWithImpl<$Res>
    implements _$ProjectCopyWith<$Res> {
  __$ProjectCopyWithImpl(this._self, this._then);

  final _Project _self;
  final $Res Function(_Project) _then;

/// Create a copy of Project
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? name = null,Object? slug = null,Object? status = null,Object? role = null,Object? dailyQuota = null,Object? monthlyQuota = null,Object? defaultLocale = null,}) {
  return _then(_Project(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,slug: null == slug ? _self.slug : slug // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,role: null == role ? _self.role : role // ignore: cast_nullable_to_non_nullable
as String,dailyQuota: null == dailyQuota ? _self.dailyQuota : dailyQuota // ignore: cast_nullable_to_non_nullable
as int,monthlyQuota: null == monthlyQuota ? _self.monthlyQuota : monthlyQuota // ignore: cast_nullable_to_non_nullable
as int,defaultLocale: null == defaultLocale ? _self.defaultLocale : defaultLocale // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}


/// @nodoc
mixin _$Message {

 String get id; String get status; String get channel; String get to; String get template; String get subject; String get body; String get priority;@JsonKey(name: 'provider_message_id') String get providerMessageId;@JsonKey(name: 'error_code') String get errorCode;@JsonKey(name: 'error_message') String get errorMessage; int get attempts;@JsonKey(name: 'is_test') bool get isTest;@JsonKey(name: 'cost_micros') int get costMicros; String get currency;@JsonKey(name: 'created_at') DateTime get createdAt;@JsonKey(name: 'sent_at') DateTime? get sentAt;@JsonKey(name: 'delivered_at') DateTime? get deliveredAt;@JsonKey(name: 'scheduled_at') DateTime? get scheduledAt; Map<String, dynamic> get metadata;
/// Create a copy of Message
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MessageCopyWith<Message> get copyWith => _$MessageCopyWithImpl<Message>(this as Message, _$identity);

  /// Serializes this Message to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Message;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Message&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.status, _this.status) || other.status == _this.status)&&(identical(other.channel, _this.channel) || other.channel == _this.channel)&&(identical(other.to, _this.to) || other.to == _this.to)&&(identical(other.template, _this.template) || other.template == _this.template)&&(identical(other.subject, _this.subject) || other.subject == _this.subject)&&(identical(other.body, _this.body) || other.body == _this.body)&&(identical(other.priority, _this.priority) || other.priority == _this.priority)&&(identical(other.providerMessageId, _this.providerMessageId) || other.providerMessageId == _this.providerMessageId)&&(identical(other.errorCode, _this.errorCode) || other.errorCode == _this.errorCode)&&(identical(other.errorMessage, _this.errorMessage) || other.errorMessage == _this.errorMessage)&&(identical(other.attempts, _this.attempts) || other.attempts == _this.attempts)&&(identical(other.isTest, _this.isTest) || other.isTest == _this.isTest)&&(identical(other.costMicros, _this.costMicros) || other.costMicros == _this.costMicros)&&(identical(other.currency, _this.currency) || other.currency == _this.currency)&&(identical(other.createdAt, _this.createdAt) || other.createdAt == _this.createdAt)&&(identical(other.sentAt, _this.sentAt) || other.sentAt == _this.sentAt)&&(identical(other.deliveredAt, _this.deliveredAt) || other.deliveredAt == _this.deliveredAt)&&(identical(other.scheduledAt, _this.scheduledAt) || other.scheduledAt == _this.scheduledAt)&&const DeepCollectionEquality().equals(other.metadata, _this.metadata));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Message;
  return Object.hashAll([runtimeType,_this.id,_this.status,_this.channel,_this.to,_this.template,_this.subject,_this.body,_this.priority,_this.providerMessageId,_this.errorCode,_this.errorMessage,_this.attempts,_this.isTest,_this.costMicros,_this.currency,_this.createdAt,_this.sentAt,_this.deliveredAt,_this.scheduledAt,const DeepCollectionEquality().hash(_this.metadata)]);
}

@override
String toString() {
  final _this = this as Message;
  return 'Message(id: ${_this.id}, status: ${_this.status}, channel: ${_this.channel}, to: ${_this.to}, template: ${_this.template}, subject: ${_this.subject}, body: ${_this.body}, priority: ${_this.priority}, providerMessageId: ${_this.providerMessageId}, errorCode: ${_this.errorCode}, errorMessage: ${_this.errorMessage}, attempts: ${_this.attempts}, isTest: ${_this.isTest}, costMicros: ${_this.costMicros}, currency: ${_this.currency}, createdAt: ${_this.createdAt}, sentAt: ${_this.sentAt}, deliveredAt: ${_this.deliveredAt}, scheduledAt: ${_this.scheduledAt}, metadata: ${_this.metadata})';
}


}

/// @nodoc
abstract mixin class $MessageCopyWith<$Res>  {
  factory $MessageCopyWith(Message value, $Res Function(Message) _then) = _$MessageCopyWithImpl;
@useResult
$Res call({
 String id, String status, String channel, String to, String template, String subject, String body, String priority,@JsonKey(name: 'provider_message_id') String providerMessageId,@JsonKey(name: 'error_code') String errorCode,@JsonKey(name: 'error_message') String errorMessage, int attempts,@JsonKey(name: 'is_test') bool isTest,@JsonKey(name: 'cost_micros') int costMicros, String currency,@JsonKey(name: 'created_at') DateTime createdAt,@JsonKey(name: 'sent_at') DateTime? sentAt,@JsonKey(name: 'delivered_at') DateTime? deliveredAt,@JsonKey(name: 'scheduled_at') DateTime? scheduledAt, Map<String, dynamic> metadata
});




}
/// @nodoc
class _$MessageCopyWithImpl<$Res>
    implements $MessageCopyWith<$Res> {
  _$MessageCopyWithImpl(this._self, this._then);

  final Message _self;
  final $Res Function(Message) _then;

/// Create a copy of Message
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? status = null,Object? channel = null,Object? to = null,Object? template = null,Object? subject = null,Object? body = null,Object? priority = null,Object? providerMessageId = null,Object? errorCode = null,Object? errorMessage = null,Object? attempts = null,Object? isTest = null,Object? costMicros = null,Object? currency = null,Object? createdAt = null,Object? sentAt = freezed,Object? deliveredAt = freezed,Object? scheduledAt = freezed,Object? metadata = null,}) {
  return _then(Message(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,to: null == to ? _self.to : to // ignore: cast_nullable_to_non_nullable
as String,template: null == template ? _self.template : template // ignore: cast_nullable_to_non_nullable
as String,subject: null == subject ? _self.subject : subject // ignore: cast_nullable_to_non_nullable
as String,body: null == body ? _self.body : body // ignore: cast_nullable_to_non_nullable
as String,priority: null == priority ? _self.priority : priority // ignore: cast_nullable_to_non_nullable
as String,providerMessageId: null == providerMessageId ? _self.providerMessageId : providerMessageId // ignore: cast_nullable_to_non_nullable
as String,errorCode: null == errorCode ? _self.errorCode : errorCode // ignore: cast_nullable_to_non_nullable
as String,errorMessage: null == errorMessage ? _self.errorMessage : errorMessage // ignore: cast_nullable_to_non_nullable
as String,attempts: null == attempts ? _self.attempts : attempts // ignore: cast_nullable_to_non_nullable
as int,isTest: null == isTest ? _self.isTest : isTest // ignore: cast_nullable_to_non_nullable
as bool,costMicros: null == costMicros ? _self.costMicros : costMicros // ignore: cast_nullable_to_non_nullable
as int,currency: null == currency ? _self.currency : currency // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,sentAt: freezed == sentAt ? _self.sentAt : sentAt // ignore: cast_nullable_to_non_nullable
as DateTime?,deliveredAt: freezed == deliveredAt ? _self.deliveredAt : deliveredAt // ignore: cast_nullable_to_non_nullable
as DateTime?,scheduledAt: freezed == scheduledAt ? _self.scheduledAt : scheduledAt // ignore: cast_nullable_to_non_nullable
as DateTime?,metadata: null == metadata ? _self.metadata : metadata // ignore: cast_nullable_to_non_nullable
as Map<String, dynamic>,
  ));
}

}


/// Adds pattern-matching-related methods to [Message].
extension MessagePatterns on Message {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Message value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Message() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Message value)  $default,){
final _that = this;
switch (_that) {
case _Message():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Message value)?  $default,){
final _that = this;
switch (_that) {
case _Message() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String status,  String channel,  String to,  String template,  String subject,  String body,  String priority, @JsonKey(name: 'provider_message_id')  String providerMessageId, @JsonKey(name: 'error_code')  String errorCode, @JsonKey(name: 'error_message')  String errorMessage,  int attempts, @JsonKey(name: 'is_test')  bool isTest, @JsonKey(name: 'cost_micros')  int costMicros,  String currency, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'sent_at')  DateTime? sentAt, @JsonKey(name: 'delivered_at')  DateTime? deliveredAt, @JsonKey(name: 'scheduled_at')  DateTime? scheduledAt,  Map<String, dynamic> metadata)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Message() when $default != null:
return $default(_that.id,_that.status,_that.channel,_that.to,_that.template,_that.subject,_that.body,_that.priority,_that.providerMessageId,_that.errorCode,_that.errorMessage,_that.attempts,_that.isTest,_that.costMicros,_that.currency,_that.createdAt,_that.sentAt,_that.deliveredAt,_that.scheduledAt,_that.metadata);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String status,  String channel,  String to,  String template,  String subject,  String body,  String priority, @JsonKey(name: 'provider_message_id')  String providerMessageId, @JsonKey(name: 'error_code')  String errorCode, @JsonKey(name: 'error_message')  String errorMessage,  int attempts, @JsonKey(name: 'is_test')  bool isTest, @JsonKey(name: 'cost_micros')  int costMicros,  String currency, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'sent_at')  DateTime? sentAt, @JsonKey(name: 'delivered_at')  DateTime? deliveredAt, @JsonKey(name: 'scheduled_at')  DateTime? scheduledAt,  Map<String, dynamic> metadata)  $default,) {final _that = this;
switch (_that) {
case _Message():
return $default(_that.id,_that.status,_that.channel,_that.to,_that.template,_that.subject,_that.body,_that.priority,_that.providerMessageId,_that.errorCode,_that.errorMessage,_that.attempts,_that.isTest,_that.costMicros,_that.currency,_that.createdAt,_that.sentAt,_that.deliveredAt,_that.scheduledAt,_that.metadata);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String status,  String channel,  String to,  String template,  String subject,  String body,  String priority, @JsonKey(name: 'provider_message_id')  String providerMessageId, @JsonKey(name: 'error_code')  String errorCode, @JsonKey(name: 'error_message')  String errorMessage,  int attempts, @JsonKey(name: 'is_test')  bool isTest, @JsonKey(name: 'cost_micros')  int costMicros,  String currency, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'sent_at')  DateTime? sentAt, @JsonKey(name: 'delivered_at')  DateTime? deliveredAt, @JsonKey(name: 'scheduled_at')  DateTime? scheduledAt,  Map<String, dynamic> metadata)?  $default,) {final _that = this;
switch (_that) {
case _Message() when $default != null:
return $default(_that.id,_that.status,_that.channel,_that.to,_that.template,_that.subject,_that.body,_that.priority,_that.providerMessageId,_that.errorCode,_that.errorMessage,_that.attempts,_that.isTest,_that.costMicros,_that.currency,_that.createdAt,_that.sentAt,_that.deliveredAt,_that.scheduledAt,_that.metadata);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Message implements Message {
  const _Message({required this.id, required this.status, required this.channel, required this.to, this.template = '', this.subject = '', this.body = '', this.priority = 'normal', @JsonKey(name: 'provider_message_id') this.providerMessageId = '', @JsonKey(name: 'error_code') this.errorCode = '', @JsonKey(name: 'error_message') this.errorMessage = '', this.attempts = 0, @JsonKey(name: 'is_test') this.isTest = false, @JsonKey(name: 'cost_micros') this.costMicros = 0, this.currency = 'TMT', @JsonKey(name: 'created_at') required this.createdAt, @JsonKey(name: 'sent_at') this.sentAt, @JsonKey(name: 'delivered_at') this.deliveredAt, @JsonKey(name: 'scheduled_at') this.scheduledAt,  Map<String, dynamic> metadata = const <String, dynamic>{}}): _metadata = metadata;
  factory _Message.fromJson(Map<String, dynamic> json) => _$MessageFromJson(json);

@override final  String id;
@override final  String status;
@override final  String channel;
@override final  String to;
@override@JsonKey() final  String template;
@override@JsonKey() final  String subject;
@override@JsonKey() final  String body;
@override@JsonKey() final  String priority;
@override@JsonKey(name: 'provider_message_id') final  String providerMessageId;
@override@JsonKey(name: 'error_code') final  String errorCode;
@override@JsonKey(name: 'error_message') final  String errorMessage;
@override@JsonKey() final  int attempts;
@override@JsonKey(name: 'is_test') final  bool isTest;
@override@JsonKey(name: 'cost_micros') final  int costMicros;
@override@JsonKey() final  String currency;
@override@JsonKey(name: 'created_at') final  DateTime createdAt;
@override@JsonKey(name: 'sent_at') final  DateTime? sentAt;
@override@JsonKey(name: 'delivered_at') final  DateTime? deliveredAt;
@override@JsonKey(name: 'scheduled_at') final  DateTime? scheduledAt;
 final  Map<String, dynamic> _metadata;
@override@JsonKey() Map<String, dynamic> get metadata {
  if (_metadata is EqualUnmodifiableMapView) return _metadata;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableMapView(_metadata);
}


/// Create a copy of Message
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MessageCopyWith<_Message> get copyWith => __$MessageCopyWithImpl<_Message>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MessageToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Message&&(identical(other.id, id) || other.id == id)&&(identical(other.status, status) || other.status == status)&&(identical(other.channel, channel) || other.channel == channel)&&(identical(other.to, to) || other.to == to)&&(identical(other.template, template) || other.template == template)&&(identical(other.subject, subject) || other.subject == subject)&&(identical(other.body, body) || other.body == body)&&(identical(other.priority, priority) || other.priority == priority)&&(identical(other.providerMessageId, providerMessageId) || other.providerMessageId == providerMessageId)&&(identical(other.errorCode, errorCode) || other.errorCode == errorCode)&&(identical(other.errorMessage, errorMessage) || other.errorMessage == errorMessage)&&(identical(other.attempts, attempts) || other.attempts == attempts)&&(identical(other.isTest, isTest) || other.isTest == isTest)&&(identical(other.costMicros, costMicros) || other.costMicros == costMicros)&&(identical(other.currency, currency) || other.currency == currency)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt)&&(identical(other.sentAt, sentAt) || other.sentAt == sentAt)&&(identical(other.deliveredAt, deliveredAt) || other.deliveredAt == deliveredAt)&&(identical(other.scheduledAt, scheduledAt) || other.scheduledAt == scheduledAt)&&const DeepCollectionEquality().equals(other.metadata, _metadata));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hashAll([runtimeType,id,status,channel,to,template,subject,body,priority,providerMessageId,errorCode,errorMessage,attempts,isTest,costMicros,currency,createdAt,sentAt,deliveredAt,scheduledAt,const DeepCollectionEquality().hash(_metadata)]);
}

@override
String toString() {
    return 'Message(id: $id, status: $status, channel: $channel, to: $to, template: $template, subject: $subject, body: $body, priority: $priority, providerMessageId: $providerMessageId, errorCode: $errorCode, errorMessage: $errorMessage, attempts: $attempts, isTest: $isTest, costMicros: $costMicros, currency: $currency, createdAt: $createdAt, sentAt: $sentAt, deliveredAt: $deliveredAt, scheduledAt: $scheduledAt, metadata: $metadata)';
}


}

/// @nodoc
abstract mixin class _$MessageCopyWith<$Res> implements $MessageCopyWith<$Res> {
  factory _$MessageCopyWith(_Message value, $Res Function(_Message) _then) = __$MessageCopyWithImpl;
@override @useResult
$Res call({
 String id, String status, String channel, String to, String template, String subject, String body, String priority,@JsonKey(name: 'provider_message_id') String providerMessageId,@JsonKey(name: 'error_code') String errorCode,@JsonKey(name: 'error_message') String errorMessage, int attempts,@JsonKey(name: 'is_test') bool isTest,@JsonKey(name: 'cost_micros') int costMicros, String currency,@JsonKey(name: 'created_at') DateTime createdAt,@JsonKey(name: 'sent_at') DateTime? sentAt,@JsonKey(name: 'delivered_at') DateTime? deliveredAt,@JsonKey(name: 'scheduled_at') DateTime? scheduledAt, Map<String, dynamic> metadata
});




}
/// @nodoc
class __$MessageCopyWithImpl<$Res>
    implements _$MessageCopyWith<$Res> {
  __$MessageCopyWithImpl(this._self, this._then);

  final _Message _self;
  final $Res Function(_Message) _then;

/// Create a copy of Message
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? status = null,Object? channel = null,Object? to = null,Object? template = null,Object? subject = null,Object? body = null,Object? priority = null,Object? providerMessageId = null,Object? errorCode = null,Object? errorMessage = null,Object? attempts = null,Object? isTest = null,Object? costMicros = null,Object? currency = null,Object? createdAt = null,Object? sentAt = freezed,Object? deliveredAt = freezed,Object? scheduledAt = freezed,Object? metadata = null,}) {
  return _then(_Message(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,to: null == to ? _self.to : to // ignore: cast_nullable_to_non_nullable
as String,template: null == template ? _self.template : template // ignore: cast_nullable_to_non_nullable
as String,subject: null == subject ? _self.subject : subject // ignore: cast_nullable_to_non_nullable
as String,body: null == body ? _self.body : body // ignore: cast_nullable_to_non_nullable
as String,priority: null == priority ? _self.priority : priority // ignore: cast_nullable_to_non_nullable
as String,providerMessageId: null == providerMessageId ? _self.providerMessageId : providerMessageId // ignore: cast_nullable_to_non_nullable
as String,errorCode: null == errorCode ? _self.errorCode : errorCode // ignore: cast_nullable_to_non_nullable
as String,errorMessage: null == errorMessage ? _self.errorMessage : errorMessage // ignore: cast_nullable_to_non_nullable
as String,attempts: null == attempts ? _self.attempts : attempts // ignore: cast_nullable_to_non_nullable
as int,isTest: null == isTest ? _self.isTest : isTest // ignore: cast_nullable_to_non_nullable
as bool,costMicros: null == costMicros ? _self.costMicros : costMicros // ignore: cast_nullable_to_non_nullable
as int,currency: null == currency ? _self.currency : currency // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,sentAt: freezed == sentAt ? _self.sentAt : sentAt // ignore: cast_nullable_to_non_nullable
as DateTime?,deliveredAt: freezed == deliveredAt ? _self.deliveredAt : deliveredAt // ignore: cast_nullable_to_non_nullable
as DateTime?,scheduledAt: freezed == scheduledAt ? _self.scheduledAt : scheduledAt // ignore: cast_nullable_to_non_nullable
as DateTime?,metadata: null == metadata ? _self._metadata : metadata // ignore: cast_nullable_to_non_nullable
as Map<String, dynamic>,
  ));
}


}


/// @nodoc
mixin _$MessageEvent {

 String get id; String get type; Map<String, dynamic> get payload;@JsonKey(name: 'created_at') DateTime get createdAt;
/// Create a copy of MessageEvent
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MessageEventCopyWith<MessageEvent> get copyWith => _$MessageEventCopyWithImpl<MessageEvent>(this as MessageEvent, _$identity);

  /// Serializes this MessageEvent to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as MessageEvent;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MessageEvent&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.type, _this.type) || other.type == _this.type)&&const DeepCollectionEquality().equals(other.payload, _this.payload)&&(identical(other.createdAt, _this.createdAt) || other.createdAt == _this.createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as MessageEvent;
  return Object.hash(runtimeType,_this.id,_this.type,const DeepCollectionEquality().hash(_this.payload),_this.createdAt);
}

@override
String toString() {
  final _this = this as MessageEvent;
  return 'MessageEvent(id: ${_this.id}, type: ${_this.type}, payload: ${_this.payload}, createdAt: ${_this.createdAt})';
}


}

/// @nodoc
abstract mixin class $MessageEventCopyWith<$Res>  {
  factory $MessageEventCopyWith(MessageEvent value, $Res Function(MessageEvent) _then) = _$MessageEventCopyWithImpl;
@useResult
$Res call({
 String id, String type, Map<String, dynamic> payload,@JsonKey(name: 'created_at') DateTime createdAt
});




}
/// @nodoc
class _$MessageEventCopyWithImpl<$Res>
    implements $MessageEventCopyWith<$Res> {
  _$MessageEventCopyWithImpl(this._self, this._then);

  final MessageEvent _self;
  final $Res Function(MessageEvent) _then;

/// Create a copy of MessageEvent
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? type = null,Object? payload = null,Object? createdAt = null,}) {
  return _then(MessageEvent(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,payload: null == payload ? _self.payload : payload // ignore: cast_nullable_to_non_nullable
as Map<String, dynamic>,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [MessageEvent].
extension MessageEventPatterns on MessageEvent {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MessageEvent value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MessageEvent() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MessageEvent value)  $default,){
final _that = this;
switch (_that) {
case _MessageEvent():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MessageEvent value)?  $default,){
final _that = this;
switch (_that) {
case _MessageEvent() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String type,  Map<String, dynamic> payload, @JsonKey(name: 'created_at')  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MessageEvent() when $default != null:
return $default(_that.id,_that.type,_that.payload,_that.createdAt);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String type,  Map<String, dynamic> payload, @JsonKey(name: 'created_at')  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _MessageEvent():
return $default(_that.id,_that.type,_that.payload,_that.createdAt);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String type,  Map<String, dynamic> payload, @JsonKey(name: 'created_at')  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _MessageEvent() when $default != null:
return $default(_that.id,_that.type,_that.payload,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _MessageEvent implements MessageEvent {
  const _MessageEvent({required this.id, required this.type,  Map<String, dynamic> payload = const <String, dynamic>{}, @JsonKey(name: 'created_at') required this.createdAt}): _payload = payload;
  factory _MessageEvent.fromJson(Map<String, dynamic> json) => _$MessageEventFromJson(json);

@override final  String id;
@override final  String type;
 final  Map<String, dynamic> _payload;
@override@JsonKey() Map<String, dynamic> get payload {
  if (_payload is EqualUnmodifiableMapView) return _payload;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableMapView(_payload);
}

@override@JsonKey(name: 'created_at') final  DateTime createdAt;

/// Create a copy of MessageEvent
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MessageEventCopyWith<_MessageEvent> get copyWith => __$MessageEventCopyWithImpl<_MessageEvent>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MessageEventToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _MessageEvent&&(identical(other.id, id) || other.id == id)&&(identical(other.type, type) || other.type == type)&&const DeepCollectionEquality().equals(other.payload, _payload)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,type,const DeepCollectionEquality().hash(_payload),createdAt);
}

@override
String toString() {
    return 'MessageEvent(id: $id, type: $type, payload: $payload, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$MessageEventCopyWith<$Res> implements $MessageEventCopyWith<$Res> {
  factory _$MessageEventCopyWith(_MessageEvent value, $Res Function(_MessageEvent) _then) = __$MessageEventCopyWithImpl;
@override @useResult
$Res call({
 String id, String type, Map<String, dynamic> payload,@JsonKey(name: 'created_at') DateTime createdAt
});




}
/// @nodoc
class __$MessageEventCopyWithImpl<$Res>
    implements _$MessageEventCopyWith<$Res> {
  __$MessageEventCopyWithImpl(this._self, this._then);

  final _MessageEvent _self;
  final $Res Function(_MessageEvent) _then;

/// Create a copy of MessageEvent
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? type = null,Object? payload = null,Object? createdAt = null,}) {
  return _then(_MessageEvent(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,payload: null == payload ? _self._payload : payload // ignore: cast_nullable_to_non_nullable
as Map<String, dynamic>,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}


/// @nodoc
mixin _$MessageDetail {

 Message get message; List<MessageEvent> get events;
/// Create a copy of MessageDetail
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MessageDetailCopyWith<MessageDetail> get copyWith => _$MessageDetailCopyWithImpl<MessageDetail>(this as MessageDetail, _$identity);

  /// Serializes this MessageDetail to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as MessageDetail;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MessageDetail&&(identical(other.message, _this.message) || other.message == _this.message)&&const DeepCollectionEquality().equals(other.events, _this.events));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as MessageDetail;
  return Object.hash(runtimeType,_this.message,const DeepCollectionEquality().hash(_this.events));
}

@override
String toString() {
  final _this = this as MessageDetail;
  return 'MessageDetail(message: ${_this.message}, events: ${_this.events})';
}


}

/// @nodoc
abstract mixin class $MessageDetailCopyWith<$Res>  {
  factory $MessageDetailCopyWith(MessageDetail value, $Res Function(MessageDetail) _then) = _$MessageDetailCopyWithImpl;
@useResult
$Res call({
 Message message, List<MessageEvent> events
});


$MessageCopyWith<$Res> get message;

}
/// @nodoc
class _$MessageDetailCopyWithImpl<$Res>
    implements $MessageDetailCopyWith<$Res> {
  _$MessageDetailCopyWithImpl(this._self, this._then);

  final MessageDetail _self;
  final $Res Function(MessageDetail) _then;

/// Create a copy of MessageDetail
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? message = null,Object? events = null,}) {
  return _then(MessageDetail(
message: null == message ? _self.message : message // ignore: cast_nullable_to_non_nullable
as Message,events: null == events ? _self.events : events // ignore: cast_nullable_to_non_nullable
as List<MessageEvent>,
  ));
}
/// Create a copy of MessageDetail
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MessageCopyWith<$Res> get message {
  
  return $MessageCopyWith<$Res>(_self.message, (value) {
    return _then(_self.copyWith(message: value));
  });
}
}


/// Adds pattern-matching-related methods to [MessageDetail].
extension MessageDetailPatterns on MessageDetail {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MessageDetail value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MessageDetail() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MessageDetail value)  $default,){
final _that = this;
switch (_that) {
case _MessageDetail():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MessageDetail value)?  $default,){
final _that = this;
switch (_that) {
case _MessageDetail() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( Message message,  List<MessageEvent> events)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MessageDetail() when $default != null:
return $default(_that.message,_that.events);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( Message message,  List<MessageEvent> events)  $default,) {final _that = this;
switch (_that) {
case _MessageDetail():
return $default(_that.message,_that.events);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( Message message,  List<MessageEvent> events)?  $default,) {final _that = this;
switch (_that) {
case _MessageDetail() when $default != null:
return $default(_that.message,_that.events);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _MessageDetail implements MessageDetail {
  const _MessageDetail({required this.message,  List<MessageEvent> events = const <MessageEvent>[]}): _events = events;
  factory _MessageDetail.fromJson(Map<String, dynamic> json) => _$MessageDetailFromJson(json);

@override final  Message message;
 final  List<MessageEvent> _events;
@override@JsonKey() List<MessageEvent> get events {
  if (_events is EqualUnmodifiableListView) return _events;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_events);
}


/// Create a copy of MessageDetail
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MessageDetailCopyWith<_MessageDetail> get copyWith => __$MessageDetailCopyWithImpl<_MessageDetail>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MessageDetailToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _MessageDetail&&(identical(other.message, message) || other.message == message)&&const DeepCollectionEquality().equals(other.events, _events));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,message,const DeepCollectionEquality().hash(_events));
}

@override
String toString() {
    return 'MessageDetail(message: $message, events: $events)';
}


}

/// @nodoc
abstract mixin class _$MessageDetailCopyWith<$Res> implements $MessageDetailCopyWith<$Res> {
  factory _$MessageDetailCopyWith(_MessageDetail value, $Res Function(_MessageDetail) _then) = __$MessageDetailCopyWithImpl;
@override @useResult
$Res call({
 Message message, List<MessageEvent> events
});


@override $MessageCopyWith<$Res> get message;

}
/// @nodoc
class __$MessageDetailCopyWithImpl<$Res>
    implements _$MessageDetailCopyWith<$Res> {
  __$MessageDetailCopyWithImpl(this._self, this._then);

  final _MessageDetail _self;
  final $Res Function(_MessageDetail) _then;

/// Create a copy of MessageDetail
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? message = null,Object? events = null,}) {
  return _then(_MessageDetail(
message: null == message ? _self.message : message // ignore: cast_nullable_to_non_nullable
as Message,events: null == events ? _self._events : events // ignore: cast_nullable_to_non_nullable
as List<MessageEvent>,
  ));
}

/// Create a copy of MessageDetail
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MessageCopyWith<$Res> get message {
  
  return $MessageCopyWith<$Res>(_self.message, (value) {
    return _then(_self.copyWith(message: value));
  });
}
}


/// @nodoc
mixin _$Template {

 String get id; String get key; String get channel; String get locale; String get subject; String get body;@JsonKey(name: 'required_vars') List<String> get requiredVars; int get version;@JsonKey(name: 'is_active') bool get isActive;
/// Create a copy of Template
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TemplateCopyWith<Template> get copyWith => _$TemplateCopyWithImpl<Template>(this as Template, _$identity);

  /// Serializes this Template to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Template;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Template&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.key, _this.key) || other.key == _this.key)&&(identical(other.channel, _this.channel) || other.channel == _this.channel)&&(identical(other.locale, _this.locale) || other.locale == _this.locale)&&(identical(other.subject, _this.subject) || other.subject == _this.subject)&&(identical(other.body, _this.body) || other.body == _this.body)&&const DeepCollectionEquality().equals(other.requiredVars, _this.requiredVars)&&(identical(other.version, _this.version) || other.version == _this.version)&&(identical(other.isActive, _this.isActive) || other.isActive == _this.isActive));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Template;
  return Object.hash(runtimeType,_this.id,_this.key,_this.channel,_this.locale,_this.subject,_this.body,const DeepCollectionEquality().hash(_this.requiredVars),_this.version,_this.isActive);
}

@override
String toString() {
  final _this = this as Template;
  return 'Template(id: ${_this.id}, key: ${_this.key}, channel: ${_this.channel}, locale: ${_this.locale}, subject: ${_this.subject}, body: ${_this.body}, requiredVars: ${_this.requiredVars}, version: ${_this.version}, isActive: ${_this.isActive})';
}


}

/// @nodoc
abstract mixin class $TemplateCopyWith<$Res>  {
  factory $TemplateCopyWith(Template value, $Res Function(Template) _then) = _$TemplateCopyWithImpl;
@useResult
$Res call({
 String id, String key, String channel, String locale, String subject, String body,@JsonKey(name: 'required_vars') List<String> requiredVars, int version,@JsonKey(name: 'is_active') bool isActive
});




}
/// @nodoc
class _$TemplateCopyWithImpl<$Res>
    implements $TemplateCopyWith<$Res> {
  _$TemplateCopyWithImpl(this._self, this._then);

  final Template _self;
  final $Res Function(Template) _then;

/// Create a copy of Template
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? key = null,Object? channel = null,Object? locale = null,Object? subject = null,Object? body = null,Object? requiredVars = null,Object? version = null,Object? isActive = null,}) {
  return _then(Template(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,key: null == key ? _self.key : key // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,locale: null == locale ? _self.locale : locale // ignore: cast_nullable_to_non_nullable
as String,subject: null == subject ? _self.subject : subject // ignore: cast_nullable_to_non_nullable
as String,body: null == body ? _self.body : body // ignore: cast_nullable_to_non_nullable
as String,requiredVars: null == requiredVars ? _self.requiredVars : requiredVars // ignore: cast_nullable_to_non_nullable
as List<String>,version: null == version ? _self.version : version // ignore: cast_nullable_to_non_nullable
as int,isActive: null == isActive ? _self.isActive : isActive // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

}


/// Adds pattern-matching-related methods to [Template].
extension TemplatePatterns on Template {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Template value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Template() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Template value)  $default,){
final _that = this;
switch (_that) {
case _Template():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Template value)?  $default,){
final _that = this;
switch (_that) {
case _Template() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String key,  String channel,  String locale,  String subject,  String body, @JsonKey(name: 'required_vars')  List<String> requiredVars,  int version, @JsonKey(name: 'is_active')  bool isActive)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Template() when $default != null:
return $default(_that.id,_that.key,_that.channel,_that.locale,_that.subject,_that.body,_that.requiredVars,_that.version,_that.isActive);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String key,  String channel,  String locale,  String subject,  String body, @JsonKey(name: 'required_vars')  List<String> requiredVars,  int version, @JsonKey(name: 'is_active')  bool isActive)  $default,) {final _that = this;
switch (_that) {
case _Template():
return $default(_that.id,_that.key,_that.channel,_that.locale,_that.subject,_that.body,_that.requiredVars,_that.version,_that.isActive);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String key,  String channel,  String locale,  String subject,  String body, @JsonKey(name: 'required_vars')  List<String> requiredVars,  int version, @JsonKey(name: 'is_active')  bool isActive)?  $default,) {final _that = this;
switch (_that) {
case _Template() when $default != null:
return $default(_that.id,_that.key,_that.channel,_that.locale,_that.subject,_that.body,_that.requiredVars,_that.version,_that.isActive);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Template implements Template {
  const _Template({required this.id, required this.key, required this.channel, this.locale = 'tk', this.subject = '', this.body = '', @JsonKey(name: 'required_vars')  List<String> requiredVars = const <String>[], this.version = 1, @JsonKey(name: 'is_active') this.isActive = true}): _requiredVars = requiredVars;
  factory _Template.fromJson(Map<String, dynamic> json) => _$TemplateFromJson(json);

@override final  String id;
@override final  String key;
@override final  String channel;
@override@JsonKey() final  String locale;
@override@JsonKey() final  String subject;
@override@JsonKey() final  String body;
 final  List<String> _requiredVars;
@override@JsonKey(name: 'required_vars') List<String> get requiredVars {
  if (_requiredVars is EqualUnmodifiableListView) return _requiredVars;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_requiredVars);
}

@override@JsonKey() final  int version;
@override@JsonKey(name: 'is_active') final  bool isActive;

/// Create a copy of Template
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TemplateCopyWith<_Template> get copyWith => __$TemplateCopyWithImpl<_Template>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$TemplateToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Template&&(identical(other.id, id) || other.id == id)&&(identical(other.key, key) || other.key == key)&&(identical(other.channel, channel) || other.channel == channel)&&(identical(other.locale, locale) || other.locale == locale)&&(identical(other.subject, subject) || other.subject == subject)&&(identical(other.body, body) || other.body == body)&&const DeepCollectionEquality().equals(other.requiredVars, _requiredVars)&&(identical(other.version, version) || other.version == version)&&(identical(other.isActive, isActive) || other.isActive == isActive));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,key,channel,locale,subject,body,const DeepCollectionEquality().hash(_requiredVars),version,isActive);
}

@override
String toString() {
    return 'Template(id: $id, key: $key, channel: $channel, locale: $locale, subject: $subject, body: $body, requiredVars: $requiredVars, version: $version, isActive: $isActive)';
}


}

/// @nodoc
abstract mixin class _$TemplateCopyWith<$Res> implements $TemplateCopyWith<$Res> {
  factory _$TemplateCopyWith(_Template value, $Res Function(_Template) _then) = __$TemplateCopyWithImpl;
@override @useResult
$Res call({
 String id, String key, String channel, String locale, String subject, String body,@JsonKey(name: 'required_vars') List<String> requiredVars, int version,@JsonKey(name: 'is_active') bool isActive
});




}
/// @nodoc
class __$TemplateCopyWithImpl<$Res>
    implements _$TemplateCopyWith<$Res> {
  __$TemplateCopyWithImpl(this._self, this._then);

  final _Template _self;
  final $Res Function(_Template) _then;

/// Create a copy of Template
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? key = null,Object? channel = null,Object? locale = null,Object? subject = null,Object? body = null,Object? requiredVars = null,Object? version = null,Object? isActive = null,}) {
  return _then(_Template(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,key: null == key ? _self.key : key // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,locale: null == locale ? _self.locale : locale // ignore: cast_nullable_to_non_nullable
as String,subject: null == subject ? _self.subject : subject // ignore: cast_nullable_to_non_nullable
as String,body: null == body ? _self.body : body // ignore: cast_nullable_to_non_nullable
as String,requiredVars: null == requiredVars ? _self._requiredVars : requiredVars // ignore: cast_nullable_to_non_nullable
as List<String>,version: null == version ? _self.version : version // ignore: cast_nullable_to_non_nullable
as int,isActive: null == isActive ? _self.isActive : isActive // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}


}


/// @nodoc
mixin _$Preview {

 String get subject; String get body;@JsonKey(name: 'required_vars') List<String> get requiredVars;@JsonKey(name: 'missing_vars') List<String> get missingVars;
/// Create a copy of Preview
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PreviewCopyWith<Preview> get copyWith => _$PreviewCopyWithImpl<Preview>(this as Preview, _$identity);

  /// Serializes this Preview to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Preview;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Preview&&(identical(other.subject, _this.subject) || other.subject == _this.subject)&&(identical(other.body, _this.body) || other.body == _this.body)&&const DeepCollectionEquality().equals(other.requiredVars, _this.requiredVars)&&const DeepCollectionEquality().equals(other.missingVars, _this.missingVars));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Preview;
  return Object.hash(runtimeType,_this.subject,_this.body,const DeepCollectionEquality().hash(_this.requiredVars),const DeepCollectionEquality().hash(_this.missingVars));
}

@override
String toString() {
  final _this = this as Preview;
  return 'Preview(subject: ${_this.subject}, body: ${_this.body}, requiredVars: ${_this.requiredVars}, missingVars: ${_this.missingVars})';
}


}

/// @nodoc
abstract mixin class $PreviewCopyWith<$Res>  {
  factory $PreviewCopyWith(Preview value, $Res Function(Preview) _then) = _$PreviewCopyWithImpl;
@useResult
$Res call({
 String subject, String body,@JsonKey(name: 'required_vars') List<String> requiredVars,@JsonKey(name: 'missing_vars') List<String> missingVars
});




}
/// @nodoc
class _$PreviewCopyWithImpl<$Res>
    implements $PreviewCopyWith<$Res> {
  _$PreviewCopyWithImpl(this._self, this._then);

  final Preview _self;
  final $Res Function(Preview) _then;

/// Create a copy of Preview
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? subject = null,Object? body = null,Object? requiredVars = null,Object? missingVars = null,}) {
  return _then(Preview(
subject: null == subject ? _self.subject : subject // ignore: cast_nullable_to_non_nullable
as String,body: null == body ? _self.body : body // ignore: cast_nullable_to_non_nullable
as String,requiredVars: null == requiredVars ? _self.requiredVars : requiredVars // ignore: cast_nullable_to_non_nullable
as List<String>,missingVars: null == missingVars ? _self.missingVars : missingVars // ignore: cast_nullable_to_non_nullable
as List<String>,
  ));
}

}


/// Adds pattern-matching-related methods to [Preview].
extension PreviewPatterns on Preview {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Preview value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Preview() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Preview value)  $default,){
final _that = this;
switch (_that) {
case _Preview():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Preview value)?  $default,){
final _that = this;
switch (_that) {
case _Preview() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String subject,  String body, @JsonKey(name: 'required_vars')  List<String> requiredVars, @JsonKey(name: 'missing_vars')  List<String> missingVars)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Preview() when $default != null:
return $default(_that.subject,_that.body,_that.requiredVars,_that.missingVars);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String subject,  String body, @JsonKey(name: 'required_vars')  List<String> requiredVars, @JsonKey(name: 'missing_vars')  List<String> missingVars)  $default,) {final _that = this;
switch (_that) {
case _Preview():
return $default(_that.subject,_that.body,_that.requiredVars,_that.missingVars);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String subject,  String body, @JsonKey(name: 'required_vars')  List<String> requiredVars, @JsonKey(name: 'missing_vars')  List<String> missingVars)?  $default,) {final _that = this;
switch (_that) {
case _Preview() when $default != null:
return $default(_that.subject,_that.body,_that.requiredVars,_that.missingVars);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Preview implements Preview {
  const _Preview({this.subject = '', this.body = '', @JsonKey(name: 'required_vars')  List<String> requiredVars = const <String>[], @JsonKey(name: 'missing_vars')  List<String> missingVars = const <String>[]}): _requiredVars = requiredVars,_missingVars = missingVars;
  factory _Preview.fromJson(Map<String, dynamic> json) => _$PreviewFromJson(json);

@override@JsonKey() final  String subject;
@override@JsonKey() final  String body;
 final  List<String> _requiredVars;
@override@JsonKey(name: 'required_vars') List<String> get requiredVars {
  if (_requiredVars is EqualUnmodifiableListView) return _requiredVars;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_requiredVars);
}

 final  List<String> _missingVars;
@override@JsonKey(name: 'missing_vars') List<String> get missingVars {
  if (_missingVars is EqualUnmodifiableListView) return _missingVars;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_missingVars);
}


/// Create a copy of Preview
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PreviewCopyWith<_Preview> get copyWith => __$PreviewCopyWithImpl<_Preview>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$PreviewToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Preview&&(identical(other.subject, subject) || other.subject == subject)&&(identical(other.body, body) || other.body == body)&&const DeepCollectionEquality().equals(other.requiredVars, _requiredVars)&&const DeepCollectionEquality().equals(other.missingVars, _missingVars));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,subject,body,const DeepCollectionEquality().hash(_requiredVars),const DeepCollectionEquality().hash(_missingVars));
}

@override
String toString() {
    return 'Preview(subject: $subject, body: $body, requiredVars: $requiredVars, missingVars: $missingVars)';
}


}

/// @nodoc
abstract mixin class _$PreviewCopyWith<$Res> implements $PreviewCopyWith<$Res> {
  factory _$PreviewCopyWith(_Preview value, $Res Function(_Preview) _then) = __$PreviewCopyWithImpl;
@override @useResult
$Res call({
 String subject, String body,@JsonKey(name: 'required_vars') List<String> requiredVars,@JsonKey(name: 'missing_vars') List<String> missingVars
});




}
/// @nodoc
class __$PreviewCopyWithImpl<$Res>
    implements _$PreviewCopyWith<$Res> {
  __$PreviewCopyWithImpl(this._self, this._then);

  final _Preview _self;
  final $Res Function(_Preview) _then;

/// Create a copy of Preview
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? subject = null,Object? body = null,Object? requiredVars = null,Object? missingVars = null,}) {
  return _then(_Preview(
subject: null == subject ? _self.subject : subject // ignore: cast_nullable_to_non_nullable
as String,body: null == body ? _self.body : body // ignore: cast_nullable_to_non_nullable
as String,requiredVars: null == requiredVars ? _self._requiredVars : requiredVars // ignore: cast_nullable_to_non_nullable
as List<String>,missingVars: null == missingVars ? _self._missingVars : missingVars // ignore: cast_nullable_to_non_nullable
as List<String>,
  ));
}


}


/// @nodoc
mixin _$ProviderInfo {

 String get id; String get name; String get channel; String get type; int get priority;@JsonKey(name: 'is_active') bool get isActive;@JsonKey(name: 'rate_limit_per_sec') int get rateLimitPerSec;
/// Create a copy of ProviderInfo
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ProviderInfoCopyWith<ProviderInfo> get copyWith => _$ProviderInfoCopyWithImpl<ProviderInfo>(this as ProviderInfo, _$identity);

  /// Serializes this ProviderInfo to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as ProviderInfo;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ProviderInfo&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.name, _this.name) || other.name == _this.name)&&(identical(other.channel, _this.channel) || other.channel == _this.channel)&&(identical(other.type, _this.type) || other.type == _this.type)&&(identical(other.priority, _this.priority) || other.priority == _this.priority)&&(identical(other.isActive, _this.isActive) || other.isActive == _this.isActive)&&(identical(other.rateLimitPerSec, _this.rateLimitPerSec) || other.rateLimitPerSec == _this.rateLimitPerSec));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as ProviderInfo;
  return Object.hash(runtimeType,_this.id,_this.name,_this.channel,_this.type,_this.priority,_this.isActive,_this.rateLimitPerSec);
}

@override
String toString() {
  final _this = this as ProviderInfo;
  return 'ProviderInfo(id: ${_this.id}, name: ${_this.name}, channel: ${_this.channel}, type: ${_this.type}, priority: ${_this.priority}, isActive: ${_this.isActive}, rateLimitPerSec: ${_this.rateLimitPerSec})';
}


}

/// @nodoc
abstract mixin class $ProviderInfoCopyWith<$Res>  {
  factory $ProviderInfoCopyWith(ProviderInfo value, $Res Function(ProviderInfo) _then) = _$ProviderInfoCopyWithImpl;
@useResult
$Res call({
 String id, String name, String channel, String type, int priority,@JsonKey(name: 'is_active') bool isActive,@JsonKey(name: 'rate_limit_per_sec') int rateLimitPerSec
});




}
/// @nodoc
class _$ProviderInfoCopyWithImpl<$Res>
    implements $ProviderInfoCopyWith<$Res> {
  _$ProviderInfoCopyWithImpl(this._self, this._then);

  final ProviderInfo _self;
  final $Res Function(ProviderInfo) _then;

/// Create a copy of ProviderInfo
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? name = null,Object? channel = null,Object? type = null,Object? priority = null,Object? isActive = null,Object? rateLimitPerSec = null,}) {
  return _then(ProviderInfo(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,priority: null == priority ? _self.priority : priority // ignore: cast_nullable_to_non_nullable
as int,isActive: null == isActive ? _self.isActive : isActive // ignore: cast_nullable_to_non_nullable
as bool,rateLimitPerSec: null == rateLimitPerSec ? _self.rateLimitPerSec : rateLimitPerSec // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [ProviderInfo].
extension ProviderInfoPatterns on ProviderInfo {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ProviderInfo value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ProviderInfo() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ProviderInfo value)  $default,){
final _that = this;
switch (_that) {
case _ProviderInfo():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ProviderInfo value)?  $default,){
final _that = this;
switch (_that) {
case _ProviderInfo() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String name,  String channel,  String type,  int priority, @JsonKey(name: 'is_active')  bool isActive, @JsonKey(name: 'rate_limit_per_sec')  int rateLimitPerSec)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ProviderInfo() when $default != null:
return $default(_that.id,_that.name,_that.channel,_that.type,_that.priority,_that.isActive,_that.rateLimitPerSec);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String name,  String channel,  String type,  int priority, @JsonKey(name: 'is_active')  bool isActive, @JsonKey(name: 'rate_limit_per_sec')  int rateLimitPerSec)  $default,) {final _that = this;
switch (_that) {
case _ProviderInfo():
return $default(_that.id,_that.name,_that.channel,_that.type,_that.priority,_that.isActive,_that.rateLimitPerSec);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String name,  String channel,  String type,  int priority, @JsonKey(name: 'is_active')  bool isActive, @JsonKey(name: 'rate_limit_per_sec')  int rateLimitPerSec)?  $default,) {final _that = this;
switch (_that) {
case _ProviderInfo() when $default != null:
return $default(_that.id,_that.name,_that.channel,_that.type,_that.priority,_that.isActive,_that.rateLimitPerSec);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ProviderInfo implements ProviderInfo {
  const _ProviderInfo({required this.id, required this.name, required this.channel, required this.type, this.priority = 100, @JsonKey(name: 'is_active') this.isActive = true, @JsonKey(name: 'rate_limit_per_sec') this.rateLimitPerSec = 0});
  factory _ProviderInfo.fromJson(Map<String, dynamic> json) => _$ProviderInfoFromJson(json);

@override final  String id;
@override final  String name;
@override final  String channel;
@override final  String type;
@override@JsonKey() final  int priority;
@override@JsonKey(name: 'is_active') final  bool isActive;
@override@JsonKey(name: 'rate_limit_per_sec') final  int rateLimitPerSec;

/// Create a copy of ProviderInfo
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ProviderInfoCopyWith<_ProviderInfo> get copyWith => __$ProviderInfoCopyWithImpl<_ProviderInfo>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ProviderInfoToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _ProviderInfo&&(identical(other.id, id) || other.id == id)&&(identical(other.name, name) || other.name == name)&&(identical(other.channel, channel) || other.channel == channel)&&(identical(other.type, type) || other.type == type)&&(identical(other.priority, priority) || other.priority == priority)&&(identical(other.isActive, isActive) || other.isActive == isActive)&&(identical(other.rateLimitPerSec, rateLimitPerSec) || other.rateLimitPerSec == rateLimitPerSec));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,name,channel,type,priority,isActive,rateLimitPerSec);
}

@override
String toString() {
    return 'ProviderInfo(id: $id, name: $name, channel: $channel, type: $type, priority: $priority, isActive: $isActive, rateLimitPerSec: $rateLimitPerSec)';
}


}

/// @nodoc
abstract mixin class _$ProviderInfoCopyWith<$Res> implements $ProviderInfoCopyWith<$Res> {
  factory _$ProviderInfoCopyWith(_ProviderInfo value, $Res Function(_ProviderInfo) _then) = __$ProviderInfoCopyWithImpl;
@override @useResult
$Res call({
 String id, String name, String channel, String type, int priority,@JsonKey(name: 'is_active') bool isActive,@JsonKey(name: 'rate_limit_per_sec') int rateLimitPerSec
});




}
/// @nodoc
class __$ProviderInfoCopyWithImpl<$Res>
    implements _$ProviderInfoCopyWith<$Res> {
  __$ProviderInfoCopyWithImpl(this._self, this._then);

  final _ProviderInfo _self;
  final $Res Function(_ProviderInfo) _then;

/// Create a copy of ProviderInfo
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? name = null,Object? channel = null,Object? type = null,Object? priority = null,Object? isActive = null,Object? rateLimitPerSec = null,}) {
  return _then(_ProviderInfo(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,priority: null == priority ? _self.priority : priority // ignore: cast_nullable_to_non_nullable
as int,isActive: null == isActive ? _self.isActive : isActive // ignore: cast_nullable_to_non_nullable
as bool,rateLimitPerSec: null == rateLimitPerSec ? _self.rateLimitPerSec : rateLimitPerSec // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$ProviderHealth {

 String get id; String get name; String get channel; String get type; String get status;@JsonKey(name: 'ok_count') int get okCount;@JsonKey(name: 'failed_count') int get failedCount;@JsonKey(name: 'last_sent_at') DateTime? get lastSentAt;
/// Create a copy of ProviderHealth
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ProviderHealthCopyWith<ProviderHealth> get copyWith => _$ProviderHealthCopyWithImpl<ProviderHealth>(this as ProviderHealth, _$identity);

  /// Serializes this ProviderHealth to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as ProviderHealth;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ProviderHealth&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.name, _this.name) || other.name == _this.name)&&(identical(other.channel, _this.channel) || other.channel == _this.channel)&&(identical(other.type, _this.type) || other.type == _this.type)&&(identical(other.status, _this.status) || other.status == _this.status)&&(identical(other.okCount, _this.okCount) || other.okCount == _this.okCount)&&(identical(other.failedCount, _this.failedCount) || other.failedCount == _this.failedCount)&&(identical(other.lastSentAt, _this.lastSentAt) || other.lastSentAt == _this.lastSentAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as ProviderHealth;
  return Object.hash(runtimeType,_this.id,_this.name,_this.channel,_this.type,_this.status,_this.okCount,_this.failedCount,_this.lastSentAt);
}

@override
String toString() {
  final _this = this as ProviderHealth;
  return 'ProviderHealth(id: ${_this.id}, name: ${_this.name}, channel: ${_this.channel}, type: ${_this.type}, status: ${_this.status}, okCount: ${_this.okCount}, failedCount: ${_this.failedCount}, lastSentAt: ${_this.lastSentAt})';
}


}

/// @nodoc
abstract mixin class $ProviderHealthCopyWith<$Res>  {
  factory $ProviderHealthCopyWith(ProviderHealth value, $Res Function(ProviderHealth) _then) = _$ProviderHealthCopyWithImpl;
@useResult
$Res call({
 String id, String name, String channel, String type, String status,@JsonKey(name: 'ok_count') int okCount,@JsonKey(name: 'failed_count') int failedCount,@JsonKey(name: 'last_sent_at') DateTime? lastSentAt
});




}
/// @nodoc
class _$ProviderHealthCopyWithImpl<$Res>
    implements $ProviderHealthCopyWith<$Res> {
  _$ProviderHealthCopyWithImpl(this._self, this._then);

  final ProviderHealth _self;
  final $Res Function(ProviderHealth) _then;

/// Create a copy of ProviderHealth
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? name = null,Object? channel = null,Object? type = null,Object? status = null,Object? okCount = null,Object? failedCount = null,Object? lastSentAt = freezed,}) {
  return _then(ProviderHealth(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,okCount: null == okCount ? _self.okCount : okCount // ignore: cast_nullable_to_non_nullable
as int,failedCount: null == failedCount ? _self.failedCount : failedCount // ignore: cast_nullable_to_non_nullable
as int,lastSentAt: freezed == lastSentAt ? _self.lastSentAt : lastSentAt // ignore: cast_nullable_to_non_nullable
as DateTime?,
  ));
}

}


/// Adds pattern-matching-related methods to [ProviderHealth].
extension ProviderHealthPatterns on ProviderHealth {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ProviderHealth value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ProviderHealth() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ProviderHealth value)  $default,){
final _that = this;
switch (_that) {
case _ProviderHealth():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ProviderHealth value)?  $default,){
final _that = this;
switch (_that) {
case _ProviderHealth() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String name,  String channel,  String type,  String status, @JsonKey(name: 'ok_count')  int okCount, @JsonKey(name: 'failed_count')  int failedCount, @JsonKey(name: 'last_sent_at')  DateTime? lastSentAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ProviderHealth() when $default != null:
return $default(_that.id,_that.name,_that.channel,_that.type,_that.status,_that.okCount,_that.failedCount,_that.lastSentAt);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String name,  String channel,  String type,  String status, @JsonKey(name: 'ok_count')  int okCount, @JsonKey(name: 'failed_count')  int failedCount, @JsonKey(name: 'last_sent_at')  DateTime? lastSentAt)  $default,) {final _that = this;
switch (_that) {
case _ProviderHealth():
return $default(_that.id,_that.name,_that.channel,_that.type,_that.status,_that.okCount,_that.failedCount,_that.lastSentAt);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String name,  String channel,  String type,  String status, @JsonKey(name: 'ok_count')  int okCount, @JsonKey(name: 'failed_count')  int failedCount, @JsonKey(name: 'last_sent_at')  DateTime? lastSentAt)?  $default,) {final _that = this;
switch (_that) {
case _ProviderHealth() when $default != null:
return $default(_that.id,_that.name,_that.channel,_that.type,_that.status,_that.okCount,_that.failedCount,_that.lastSentAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ProviderHealth implements ProviderHealth {
  const _ProviderHealth({required this.id, required this.name, required this.channel, required this.type, this.status = 'idle', @JsonKey(name: 'ok_count') this.okCount = 0, @JsonKey(name: 'failed_count') this.failedCount = 0, @JsonKey(name: 'last_sent_at') this.lastSentAt});
  factory _ProviderHealth.fromJson(Map<String, dynamic> json) => _$ProviderHealthFromJson(json);

@override final  String id;
@override final  String name;
@override final  String channel;
@override final  String type;
@override@JsonKey() final  String status;
@override@JsonKey(name: 'ok_count') final  int okCount;
@override@JsonKey(name: 'failed_count') final  int failedCount;
@override@JsonKey(name: 'last_sent_at') final  DateTime? lastSentAt;

/// Create a copy of ProviderHealth
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ProviderHealthCopyWith<_ProviderHealth> get copyWith => __$ProviderHealthCopyWithImpl<_ProviderHealth>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ProviderHealthToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _ProviderHealth&&(identical(other.id, id) || other.id == id)&&(identical(other.name, name) || other.name == name)&&(identical(other.channel, channel) || other.channel == channel)&&(identical(other.type, type) || other.type == type)&&(identical(other.status, status) || other.status == status)&&(identical(other.okCount, okCount) || other.okCount == okCount)&&(identical(other.failedCount, failedCount) || other.failedCount == failedCount)&&(identical(other.lastSentAt, lastSentAt) || other.lastSentAt == lastSentAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,name,channel,type,status,okCount,failedCount,lastSentAt);
}

@override
String toString() {
    return 'ProviderHealth(id: $id, name: $name, channel: $channel, type: $type, status: $status, okCount: $okCount, failedCount: $failedCount, lastSentAt: $lastSentAt)';
}


}

/// @nodoc
abstract mixin class _$ProviderHealthCopyWith<$Res> implements $ProviderHealthCopyWith<$Res> {
  factory _$ProviderHealthCopyWith(_ProviderHealth value, $Res Function(_ProviderHealth) _then) = __$ProviderHealthCopyWithImpl;
@override @useResult
$Res call({
 String id, String name, String channel, String type, String status,@JsonKey(name: 'ok_count') int okCount,@JsonKey(name: 'failed_count') int failedCount,@JsonKey(name: 'last_sent_at') DateTime? lastSentAt
});




}
/// @nodoc
class __$ProviderHealthCopyWithImpl<$Res>
    implements _$ProviderHealthCopyWith<$Res> {
  __$ProviderHealthCopyWithImpl(this._self, this._then);

  final _ProviderHealth _self;
  final $Res Function(_ProviderHealth) _then;

/// Create a copy of ProviderHealth
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? name = null,Object? channel = null,Object? type = null,Object? status = null,Object? okCount = null,Object? failedCount = null,Object? lastSentAt = freezed,}) {
  return _then(_ProviderHealth(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,okCount: null == okCount ? _self.okCount : okCount // ignore: cast_nullable_to_non_nullable
as int,failedCount: null == failedCount ? _self.failedCount : failedCount // ignore: cast_nullable_to_non_nullable
as int,lastSentAt: freezed == lastSentAt ? _self.lastSentAt : lastSentAt // ignore: cast_nullable_to_non_nullable
as DateTime?,
  ));
}


}


/// @nodoc
mixin _$QueueStats {

 String get queue; int get pending; int get active; int get scheduled; int get retry;@JsonKey(name: 'processed_today') int get processedToday;@JsonKey(name: 'failed_today') int get failedToday;
/// Create a copy of QueueStats
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$QueueStatsCopyWith<QueueStats> get copyWith => _$QueueStatsCopyWithImpl<QueueStats>(this as QueueStats, _$identity);

  /// Serializes this QueueStats to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as QueueStats;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is QueueStats&&(identical(other.queue, _this.queue) || other.queue == _this.queue)&&(identical(other.pending, _this.pending) || other.pending == _this.pending)&&(identical(other.active, _this.active) || other.active == _this.active)&&(identical(other.scheduled, _this.scheduled) || other.scheduled == _this.scheduled)&&(identical(other.retry, _this.retry) || other.retry == _this.retry)&&(identical(other.processedToday, _this.processedToday) || other.processedToday == _this.processedToday)&&(identical(other.failedToday, _this.failedToday) || other.failedToday == _this.failedToday));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as QueueStats;
  return Object.hash(runtimeType,_this.queue,_this.pending,_this.active,_this.scheduled,_this.retry,_this.processedToday,_this.failedToday);
}

@override
String toString() {
  final _this = this as QueueStats;
  return 'QueueStats(queue: ${_this.queue}, pending: ${_this.pending}, active: ${_this.active}, scheduled: ${_this.scheduled}, retry: ${_this.retry}, processedToday: ${_this.processedToday}, failedToday: ${_this.failedToday})';
}


}

/// @nodoc
abstract mixin class $QueueStatsCopyWith<$Res>  {
  factory $QueueStatsCopyWith(QueueStats value, $Res Function(QueueStats) _then) = _$QueueStatsCopyWithImpl;
@useResult
$Res call({
 String queue, int pending, int active, int scheduled, int retry,@JsonKey(name: 'processed_today') int processedToday,@JsonKey(name: 'failed_today') int failedToday
});




}
/// @nodoc
class _$QueueStatsCopyWithImpl<$Res>
    implements $QueueStatsCopyWith<$Res> {
  _$QueueStatsCopyWithImpl(this._self, this._then);

  final QueueStats _self;
  final $Res Function(QueueStats) _then;

/// Create a copy of QueueStats
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? queue = null,Object? pending = null,Object? active = null,Object? scheduled = null,Object? retry = null,Object? processedToday = null,Object? failedToday = null,}) {
  return _then(QueueStats(
queue: null == queue ? _self.queue : queue // ignore: cast_nullable_to_non_nullable
as String,pending: null == pending ? _self.pending : pending // ignore: cast_nullable_to_non_nullable
as int,active: null == active ? _self.active : active // ignore: cast_nullable_to_non_nullable
as int,scheduled: null == scheduled ? _self.scheduled : scheduled // ignore: cast_nullable_to_non_nullable
as int,retry: null == retry ? _self.retry : retry // ignore: cast_nullable_to_non_nullable
as int,processedToday: null == processedToday ? _self.processedToday : processedToday // ignore: cast_nullable_to_non_nullable
as int,failedToday: null == failedToday ? _self.failedToday : failedToday // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [QueueStats].
extension QueueStatsPatterns on QueueStats {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _QueueStats value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _QueueStats() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _QueueStats value)  $default,){
final _that = this;
switch (_that) {
case _QueueStats():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _QueueStats value)?  $default,){
final _that = this;
switch (_that) {
case _QueueStats() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String queue,  int pending,  int active,  int scheduled,  int retry, @JsonKey(name: 'processed_today')  int processedToday, @JsonKey(name: 'failed_today')  int failedToday)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _QueueStats() when $default != null:
return $default(_that.queue,_that.pending,_that.active,_that.scheduled,_that.retry,_that.processedToday,_that.failedToday);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String queue,  int pending,  int active,  int scheduled,  int retry, @JsonKey(name: 'processed_today')  int processedToday, @JsonKey(name: 'failed_today')  int failedToday)  $default,) {final _that = this;
switch (_that) {
case _QueueStats():
return $default(_that.queue,_that.pending,_that.active,_that.scheduled,_that.retry,_that.processedToday,_that.failedToday);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String queue,  int pending,  int active,  int scheduled,  int retry, @JsonKey(name: 'processed_today')  int processedToday, @JsonKey(name: 'failed_today')  int failedToday)?  $default,) {final _that = this;
switch (_that) {
case _QueueStats() when $default != null:
return $default(_that.queue,_that.pending,_that.active,_that.scheduled,_that.retry,_that.processedToday,_that.failedToday);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _QueueStats implements QueueStats {
  const _QueueStats({required this.queue, this.pending = 0, this.active = 0, this.scheduled = 0, this.retry = 0, @JsonKey(name: 'processed_today') this.processedToday = 0, @JsonKey(name: 'failed_today') this.failedToday = 0});
  factory _QueueStats.fromJson(Map<String, dynamic> json) => _$QueueStatsFromJson(json);

@override final  String queue;
@override@JsonKey() final  int pending;
@override@JsonKey() final  int active;
@override@JsonKey() final  int scheduled;
@override@JsonKey() final  int retry;
@override@JsonKey(name: 'processed_today') final  int processedToday;
@override@JsonKey(name: 'failed_today') final  int failedToday;

/// Create a copy of QueueStats
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$QueueStatsCopyWith<_QueueStats> get copyWith => __$QueueStatsCopyWithImpl<_QueueStats>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$QueueStatsToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _QueueStats&&(identical(other.queue, queue) || other.queue == queue)&&(identical(other.pending, pending) || other.pending == pending)&&(identical(other.active, active) || other.active == active)&&(identical(other.scheduled, scheduled) || other.scheduled == scheduled)&&(identical(other.retry, retry) || other.retry == retry)&&(identical(other.processedToday, processedToday) || other.processedToday == processedToday)&&(identical(other.failedToday, failedToday) || other.failedToday == failedToday));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,queue,pending,active,scheduled,retry,processedToday,failedToday);
}

@override
String toString() {
    return 'QueueStats(queue: $queue, pending: $pending, active: $active, scheduled: $scheduled, retry: $retry, processedToday: $processedToday, failedToday: $failedToday)';
}


}

/// @nodoc
abstract mixin class _$QueueStatsCopyWith<$Res> implements $QueueStatsCopyWith<$Res> {
  factory _$QueueStatsCopyWith(_QueueStats value, $Res Function(_QueueStats) _then) = __$QueueStatsCopyWithImpl;
@override @useResult
$Res call({
 String queue, int pending, int active, int scheduled, int retry,@JsonKey(name: 'processed_today') int processedToday,@JsonKey(name: 'failed_today') int failedToday
});




}
/// @nodoc
class __$QueueStatsCopyWithImpl<$Res>
    implements _$QueueStatsCopyWith<$Res> {
  __$QueueStatsCopyWithImpl(this._self, this._then);

  final _QueueStats _self;
  final $Res Function(_QueueStats) _then;

/// Create a copy of QueueStats
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? queue = null,Object? pending = null,Object? active = null,Object? scheduled = null,Object? retry = null,Object? processedToday = null,Object? failedToday = null,}) {
  return _then(_QueueStats(
queue: null == queue ? _self.queue : queue // ignore: cast_nullable_to_non_nullable
as String,pending: null == pending ? _self.pending : pending // ignore: cast_nullable_to_non_nullable
as int,active: null == active ? _self.active : active // ignore: cast_nullable_to_non_nullable
as int,scheduled: null == scheduled ? _self.scheduled : scheduled // ignore: cast_nullable_to_non_nullable
as int,retry: null == retry ? _self.retry : retry // ignore: cast_nullable_to_non_nullable
as int,processedToday: null == processedToday ? _self.processedToday : processedToday // ignore: cast_nullable_to_non_nullable
as int,failedToday: null == failedToday ? _self.failedToday : failedToday // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$Health {

 List<QueueStats> get queues; List<ProviderHealth> get providers; int get contacts;
/// Create a copy of Health
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$HealthCopyWith<Health> get copyWith => _$HealthCopyWithImpl<Health>(this as Health, _$identity);

  /// Serializes this Health to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Health;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Health&&const DeepCollectionEquality().equals(other.queues, _this.queues)&&const DeepCollectionEquality().equals(other.providers, _this.providers)&&(identical(other.contacts, _this.contacts) || other.contacts == _this.contacts));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Health;
  return Object.hash(runtimeType,const DeepCollectionEquality().hash(_this.queues),const DeepCollectionEquality().hash(_this.providers),_this.contacts);
}

@override
String toString() {
  final _this = this as Health;
  return 'Health(queues: ${_this.queues}, providers: ${_this.providers}, contacts: ${_this.contacts})';
}


}

/// @nodoc
abstract mixin class $HealthCopyWith<$Res>  {
  factory $HealthCopyWith(Health value, $Res Function(Health) _then) = _$HealthCopyWithImpl;
@useResult
$Res call({
 List<QueueStats> queues, List<ProviderHealth> providers, int contacts
});




}
/// @nodoc
class _$HealthCopyWithImpl<$Res>
    implements $HealthCopyWith<$Res> {
  _$HealthCopyWithImpl(this._self, this._then);

  final Health _self;
  final $Res Function(Health) _then;

/// Create a copy of Health
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? queues = null,Object? providers = null,Object? contacts = null,}) {
  return _then(Health(
queues: null == queues ? _self.queues : queues // ignore: cast_nullable_to_non_nullable
as List<QueueStats>,providers: null == providers ? _self.providers : providers // ignore: cast_nullable_to_non_nullable
as List<ProviderHealth>,contacts: null == contacts ? _self.contacts : contacts // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [Health].
extension HealthPatterns on Health {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Health value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Health() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Health value)  $default,){
final _that = this;
switch (_that) {
case _Health():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Health value)?  $default,){
final _that = this;
switch (_that) {
case _Health() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( List<QueueStats> queues,  List<ProviderHealth> providers,  int contacts)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Health() when $default != null:
return $default(_that.queues,_that.providers,_that.contacts);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( List<QueueStats> queues,  List<ProviderHealth> providers,  int contacts)  $default,) {final _that = this;
switch (_that) {
case _Health():
return $default(_that.queues,_that.providers,_that.contacts);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( List<QueueStats> queues,  List<ProviderHealth> providers,  int contacts)?  $default,) {final _that = this;
switch (_that) {
case _Health() when $default != null:
return $default(_that.queues,_that.providers,_that.contacts);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Health implements Health {
  const _Health({ List<QueueStats> queues = const <QueueStats>[],  List<ProviderHealth> providers = const <ProviderHealth>[], this.contacts = 0}): _queues = queues,_providers = providers;
  factory _Health.fromJson(Map<String, dynamic> json) => _$HealthFromJson(json);

 final  List<QueueStats> _queues;
@override@JsonKey() List<QueueStats> get queues {
  if (_queues is EqualUnmodifiableListView) return _queues;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_queues);
}

 final  List<ProviderHealth> _providers;
@override@JsonKey() List<ProviderHealth> get providers {
  if (_providers is EqualUnmodifiableListView) return _providers;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_providers);
}

@override@JsonKey() final  int contacts;

/// Create a copy of Health
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$HealthCopyWith<_Health> get copyWith => __$HealthCopyWithImpl<_Health>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$HealthToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Health&&const DeepCollectionEquality().equals(other.queues, _queues)&&const DeepCollectionEquality().equals(other.providers, _providers)&&(identical(other.contacts, contacts) || other.contacts == contacts));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,const DeepCollectionEquality().hash(_queues),const DeepCollectionEquality().hash(_providers),contacts);
}

@override
String toString() {
    return 'Health(queues: $queues, providers: $providers, contacts: $contacts)';
}


}

/// @nodoc
abstract mixin class _$HealthCopyWith<$Res> implements $HealthCopyWith<$Res> {
  factory _$HealthCopyWith(_Health value, $Res Function(_Health) _then) = __$HealthCopyWithImpl;
@override @useResult
$Res call({
 List<QueueStats> queues, List<ProviderHealth> providers, int contacts
});




}
/// @nodoc
class __$HealthCopyWithImpl<$Res>
    implements _$HealthCopyWith<$Res> {
  __$HealthCopyWithImpl(this._self, this._then);

  final _Health _self;
  final $Res Function(_Health) _then;

/// Create a copy of Health
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? queues = null,Object? providers = null,Object? contacts = null,}) {
  return _then(_Health(
queues: null == queues ? _self._queues : queues // ignore: cast_nullable_to_non_nullable
as List<QueueStats>,providers: null == providers ? _self._providers : providers // ignore: cast_nullable_to_non_nullable
as List<ProviderHealth>,contacts: null == contacts ? _self.contacts : contacts // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$ApiKey {

 String get id; String get name; String get prefix; String get hint; String? get key; List<String> get scopes;@JsonKey(name: 'last_used_at') DateTime? get lastUsedAt;@JsonKey(name: 'revoked_at') DateTime? get revokedAt;@JsonKey(name: 'created_at') DateTime get createdAt;
/// Create a copy of ApiKey
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ApiKeyCopyWith<ApiKey> get copyWith => _$ApiKeyCopyWithImpl<ApiKey>(this as ApiKey, _$identity);

  /// Serializes this ApiKey to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as ApiKey;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ApiKey&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.name, _this.name) || other.name == _this.name)&&(identical(other.prefix, _this.prefix) || other.prefix == _this.prefix)&&(identical(other.hint, _this.hint) || other.hint == _this.hint)&&(identical(other.key, _this.key) || other.key == _this.key)&&const DeepCollectionEquality().equals(other.scopes, _this.scopes)&&(identical(other.lastUsedAt, _this.lastUsedAt) || other.lastUsedAt == _this.lastUsedAt)&&(identical(other.revokedAt, _this.revokedAt) || other.revokedAt == _this.revokedAt)&&(identical(other.createdAt, _this.createdAt) || other.createdAt == _this.createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as ApiKey;
  return Object.hash(runtimeType,_this.id,_this.name,_this.prefix,_this.hint,_this.key,const DeepCollectionEquality().hash(_this.scopes),_this.lastUsedAt,_this.revokedAt,_this.createdAt);
}

@override
String toString() {
  final _this = this as ApiKey;
  return 'ApiKey(id: ${_this.id}, name: ${_this.name}, prefix: ${_this.prefix}, hint: ${_this.hint}, key: ${_this.key}, scopes: ${_this.scopes}, lastUsedAt: ${_this.lastUsedAt}, revokedAt: ${_this.revokedAt}, createdAt: ${_this.createdAt})';
}


}

/// @nodoc
abstract mixin class $ApiKeyCopyWith<$Res>  {
  factory $ApiKeyCopyWith(ApiKey value, $Res Function(ApiKey) _then) = _$ApiKeyCopyWithImpl;
@useResult
$Res call({
 String id, String name, String prefix, String hint, String? key, List<String> scopes,@JsonKey(name: 'last_used_at') DateTime? lastUsedAt,@JsonKey(name: 'revoked_at') DateTime? revokedAt,@JsonKey(name: 'created_at') DateTime createdAt
});




}
/// @nodoc
class _$ApiKeyCopyWithImpl<$Res>
    implements $ApiKeyCopyWith<$Res> {
  _$ApiKeyCopyWithImpl(this._self, this._then);

  final ApiKey _self;
  final $Res Function(ApiKey) _then;

/// Create a copy of ApiKey
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? name = null,Object? prefix = null,Object? hint = null,Object? key = freezed,Object? scopes = null,Object? lastUsedAt = freezed,Object? revokedAt = freezed,Object? createdAt = null,}) {
  return _then(ApiKey(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,prefix: null == prefix ? _self.prefix : prefix // ignore: cast_nullable_to_non_nullable
as String,hint: null == hint ? _self.hint : hint // ignore: cast_nullable_to_non_nullable
as String,key: freezed == key ? _self.key : key // ignore: cast_nullable_to_non_nullable
as String?,scopes: null == scopes ? _self.scopes : scopes // ignore: cast_nullable_to_non_nullable
as List<String>,lastUsedAt: freezed == lastUsedAt ? _self.lastUsedAt : lastUsedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,revokedAt: freezed == revokedAt ? _self.revokedAt : revokedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [ApiKey].
extension ApiKeyPatterns on ApiKey {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ApiKey value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ApiKey() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ApiKey value)  $default,){
final _that = this;
switch (_that) {
case _ApiKey():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ApiKey value)?  $default,){
final _that = this;
switch (_that) {
case _ApiKey() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String name,  String prefix,  String hint,  String? key,  List<String> scopes, @JsonKey(name: 'last_used_at')  DateTime? lastUsedAt, @JsonKey(name: 'revoked_at')  DateTime? revokedAt, @JsonKey(name: 'created_at')  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ApiKey() when $default != null:
return $default(_that.id,_that.name,_that.prefix,_that.hint,_that.key,_that.scopes,_that.lastUsedAt,_that.revokedAt,_that.createdAt);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String name,  String prefix,  String hint,  String? key,  List<String> scopes, @JsonKey(name: 'last_used_at')  DateTime? lastUsedAt, @JsonKey(name: 'revoked_at')  DateTime? revokedAt, @JsonKey(name: 'created_at')  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _ApiKey():
return $default(_that.id,_that.name,_that.prefix,_that.hint,_that.key,_that.scopes,_that.lastUsedAt,_that.revokedAt,_that.createdAt);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String name,  String prefix,  String hint,  String? key,  List<String> scopes, @JsonKey(name: 'last_used_at')  DateTime? lastUsedAt, @JsonKey(name: 'revoked_at')  DateTime? revokedAt, @JsonKey(name: 'created_at')  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _ApiKey() when $default != null:
return $default(_that.id,_that.name,_that.prefix,_that.hint,_that.key,_that.scopes,_that.lastUsedAt,_that.revokedAt,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ApiKey implements ApiKey {
  const _ApiKey({required this.id, required this.name, required this.prefix, required this.hint, this.key,  List<String> scopes = const <String>[], @JsonKey(name: 'last_used_at') this.lastUsedAt, @JsonKey(name: 'revoked_at') this.revokedAt, @JsonKey(name: 'created_at') required this.createdAt}): _scopes = scopes;
  factory _ApiKey.fromJson(Map<String, dynamic> json) => _$ApiKeyFromJson(json);

@override final  String id;
@override final  String name;
@override final  String prefix;
@override final  String hint;
@override final  String? key;
 final  List<String> _scopes;
@override@JsonKey() List<String> get scopes {
  if (_scopes is EqualUnmodifiableListView) return _scopes;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_scopes);
}

@override@JsonKey(name: 'last_used_at') final  DateTime? lastUsedAt;
@override@JsonKey(name: 'revoked_at') final  DateTime? revokedAt;
@override@JsonKey(name: 'created_at') final  DateTime createdAt;

/// Create a copy of ApiKey
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ApiKeyCopyWith<_ApiKey> get copyWith => __$ApiKeyCopyWithImpl<_ApiKey>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ApiKeyToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _ApiKey&&(identical(other.id, id) || other.id == id)&&(identical(other.name, name) || other.name == name)&&(identical(other.prefix, prefix) || other.prefix == prefix)&&(identical(other.hint, hint) || other.hint == hint)&&(identical(other.key, key) || other.key == key)&&const DeepCollectionEquality().equals(other.scopes, _scopes)&&(identical(other.lastUsedAt, lastUsedAt) || other.lastUsedAt == lastUsedAt)&&(identical(other.revokedAt, revokedAt) || other.revokedAt == revokedAt)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,name,prefix,hint,key,const DeepCollectionEquality().hash(_scopes),lastUsedAt,revokedAt,createdAt);
}

@override
String toString() {
    return 'ApiKey(id: $id, name: $name, prefix: $prefix, hint: $hint, key: $key, scopes: $scopes, lastUsedAt: $lastUsedAt, revokedAt: $revokedAt, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$ApiKeyCopyWith<$Res> implements $ApiKeyCopyWith<$Res> {
  factory _$ApiKeyCopyWith(_ApiKey value, $Res Function(_ApiKey) _then) = __$ApiKeyCopyWithImpl;
@override @useResult
$Res call({
 String id, String name, String prefix, String hint, String? key, List<String> scopes,@JsonKey(name: 'last_used_at') DateTime? lastUsedAt,@JsonKey(name: 'revoked_at') DateTime? revokedAt,@JsonKey(name: 'created_at') DateTime createdAt
});




}
/// @nodoc
class __$ApiKeyCopyWithImpl<$Res>
    implements _$ApiKeyCopyWith<$Res> {
  __$ApiKeyCopyWithImpl(this._self, this._then);

  final _ApiKey _self;
  final $Res Function(_ApiKey) _then;

/// Create a copy of ApiKey
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? name = null,Object? prefix = null,Object? hint = null,Object? key = freezed,Object? scopes = null,Object? lastUsedAt = freezed,Object? revokedAt = freezed,Object? createdAt = null,}) {
  return _then(_ApiKey(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,prefix: null == prefix ? _self.prefix : prefix // ignore: cast_nullable_to_non_nullable
as String,hint: null == hint ? _self.hint : hint // ignore: cast_nullable_to_non_nullable
as String,key: freezed == key ? _self.key : key // ignore: cast_nullable_to_non_nullable
as String?,scopes: null == scopes ? _self._scopes : scopes // ignore: cast_nullable_to_non_nullable
as List<String>,lastUsedAt: freezed == lastUsedAt ? _self.lastUsedAt : lastUsedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,revokedAt: freezed == revokedAt ? _self.revokedAt : revokedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}


/// @nodoc
mixin _$Totals {

 int get total; int get sent; int get delivered; int get failed; int get pending;@JsonKey(name: 'cost_micros') int get costMicros;
/// Create a copy of Totals
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TotalsCopyWith<Totals> get copyWith => _$TotalsCopyWithImpl<Totals>(this as Totals, _$identity);

  /// Serializes this Totals to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Totals;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Totals&&(identical(other.total, _this.total) || other.total == _this.total)&&(identical(other.sent, _this.sent) || other.sent == _this.sent)&&(identical(other.delivered, _this.delivered) || other.delivered == _this.delivered)&&(identical(other.failed, _this.failed) || other.failed == _this.failed)&&(identical(other.pending, _this.pending) || other.pending == _this.pending)&&(identical(other.costMicros, _this.costMicros) || other.costMicros == _this.costMicros));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Totals;
  return Object.hash(runtimeType,_this.total,_this.sent,_this.delivered,_this.failed,_this.pending,_this.costMicros);
}

@override
String toString() {
  final _this = this as Totals;
  return 'Totals(total: ${_this.total}, sent: ${_this.sent}, delivered: ${_this.delivered}, failed: ${_this.failed}, pending: ${_this.pending}, costMicros: ${_this.costMicros})';
}


}

/// @nodoc
abstract mixin class $TotalsCopyWith<$Res>  {
  factory $TotalsCopyWith(Totals value, $Res Function(Totals) _then) = _$TotalsCopyWithImpl;
@useResult
$Res call({
 int total, int sent, int delivered, int failed, int pending,@JsonKey(name: 'cost_micros') int costMicros
});




}
/// @nodoc
class _$TotalsCopyWithImpl<$Res>
    implements $TotalsCopyWith<$Res> {
  _$TotalsCopyWithImpl(this._self, this._then);

  final Totals _self;
  final $Res Function(Totals) _then;

/// Create a copy of Totals
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? total = null,Object? sent = null,Object? delivered = null,Object? failed = null,Object? pending = null,Object? costMicros = null,}) {
  return _then(Totals(
total: null == total ? _self.total : total // ignore: cast_nullable_to_non_nullable
as int,sent: null == sent ? _self.sent : sent // ignore: cast_nullable_to_non_nullable
as int,delivered: null == delivered ? _self.delivered : delivered // ignore: cast_nullable_to_non_nullable
as int,failed: null == failed ? _self.failed : failed // ignore: cast_nullable_to_non_nullable
as int,pending: null == pending ? _self.pending : pending // ignore: cast_nullable_to_non_nullable
as int,costMicros: null == costMicros ? _self.costMicros : costMicros // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [Totals].
extension TotalsPatterns on Totals {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Totals value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Totals() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Totals value)  $default,){
final _that = this;
switch (_that) {
case _Totals():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Totals value)?  $default,){
final _that = this;
switch (_that) {
case _Totals() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( int total,  int sent,  int delivered,  int failed,  int pending, @JsonKey(name: 'cost_micros')  int costMicros)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Totals() when $default != null:
return $default(_that.total,_that.sent,_that.delivered,_that.failed,_that.pending,_that.costMicros);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( int total,  int sent,  int delivered,  int failed,  int pending, @JsonKey(name: 'cost_micros')  int costMicros)  $default,) {final _that = this;
switch (_that) {
case _Totals():
return $default(_that.total,_that.sent,_that.delivered,_that.failed,_that.pending,_that.costMicros);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( int total,  int sent,  int delivered,  int failed,  int pending, @JsonKey(name: 'cost_micros')  int costMicros)?  $default,) {final _that = this;
switch (_that) {
case _Totals() when $default != null:
return $default(_that.total,_that.sent,_that.delivered,_that.failed,_that.pending,_that.costMicros);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Totals implements Totals {
  const _Totals({this.total = 0, this.sent = 0, this.delivered = 0, this.failed = 0, this.pending = 0, @JsonKey(name: 'cost_micros') this.costMicros = 0});
  factory _Totals.fromJson(Map<String, dynamic> json) => _$TotalsFromJson(json);

@override@JsonKey() final  int total;
@override@JsonKey() final  int sent;
@override@JsonKey() final  int delivered;
@override@JsonKey() final  int failed;
@override@JsonKey() final  int pending;
@override@JsonKey(name: 'cost_micros') final  int costMicros;

/// Create a copy of Totals
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TotalsCopyWith<_Totals> get copyWith => __$TotalsCopyWithImpl<_Totals>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$TotalsToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Totals&&(identical(other.total, total) || other.total == total)&&(identical(other.sent, sent) || other.sent == sent)&&(identical(other.delivered, delivered) || other.delivered == delivered)&&(identical(other.failed, failed) || other.failed == failed)&&(identical(other.pending, pending) || other.pending == pending)&&(identical(other.costMicros, costMicros) || other.costMicros == costMicros));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,total,sent,delivered,failed,pending,costMicros);
}

@override
String toString() {
    return 'Totals(total: $total, sent: $sent, delivered: $delivered, failed: $failed, pending: $pending, costMicros: $costMicros)';
}


}

/// @nodoc
abstract mixin class _$TotalsCopyWith<$Res> implements $TotalsCopyWith<$Res> {
  factory _$TotalsCopyWith(_Totals value, $Res Function(_Totals) _then) = __$TotalsCopyWithImpl;
@override @useResult
$Res call({
 int total, int sent, int delivered, int failed, int pending,@JsonKey(name: 'cost_micros') int costMicros
});




}
/// @nodoc
class __$TotalsCopyWithImpl<$Res>
    implements _$TotalsCopyWith<$Res> {
  __$TotalsCopyWithImpl(this._self, this._then);

  final _Totals _self;
  final $Res Function(_Totals) _then;

/// Create a copy of Totals
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? total = null,Object? sent = null,Object? delivered = null,Object? failed = null,Object? pending = null,Object? costMicros = null,}) {
  return _then(_Totals(
total: null == total ? _self.total : total // ignore: cast_nullable_to_non_nullable
as int,sent: null == sent ? _self.sent : sent // ignore: cast_nullable_to_non_nullable
as int,delivered: null == delivered ? _self.delivered : delivered // ignore: cast_nullable_to_non_nullable
as int,failed: null == failed ? _self.failed : failed // ignore: cast_nullable_to_non_nullable
as int,pending: null == pending ? _self.pending : pending // ignore: cast_nullable_to_non_nullable
as int,costMicros: null == costMicros ? _self.costMicros : costMicros // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$DailyPoint {

 String get day; String get channel; int get total; int get sent; int get delivered; int get failed;
/// Create a copy of DailyPoint
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$DailyPointCopyWith<DailyPoint> get copyWith => _$DailyPointCopyWithImpl<DailyPoint>(this as DailyPoint, _$identity);

  /// Serializes this DailyPoint to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as DailyPoint;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is DailyPoint&&(identical(other.day, _this.day) || other.day == _this.day)&&(identical(other.channel, _this.channel) || other.channel == _this.channel)&&(identical(other.total, _this.total) || other.total == _this.total)&&(identical(other.sent, _this.sent) || other.sent == _this.sent)&&(identical(other.delivered, _this.delivered) || other.delivered == _this.delivered)&&(identical(other.failed, _this.failed) || other.failed == _this.failed));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as DailyPoint;
  return Object.hash(runtimeType,_this.day,_this.channel,_this.total,_this.sent,_this.delivered,_this.failed);
}

@override
String toString() {
  final _this = this as DailyPoint;
  return 'DailyPoint(day: ${_this.day}, channel: ${_this.channel}, total: ${_this.total}, sent: ${_this.sent}, delivered: ${_this.delivered}, failed: ${_this.failed})';
}


}

/// @nodoc
abstract mixin class $DailyPointCopyWith<$Res>  {
  factory $DailyPointCopyWith(DailyPoint value, $Res Function(DailyPoint) _then) = _$DailyPointCopyWithImpl;
@useResult
$Res call({
 String day, String channel, int total, int sent, int delivered, int failed
});




}
/// @nodoc
class _$DailyPointCopyWithImpl<$Res>
    implements $DailyPointCopyWith<$Res> {
  _$DailyPointCopyWithImpl(this._self, this._then);

  final DailyPoint _self;
  final $Res Function(DailyPoint) _then;

/// Create a copy of DailyPoint
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? day = null,Object? channel = null,Object? total = null,Object? sent = null,Object? delivered = null,Object? failed = null,}) {
  return _then(DailyPoint(
day: null == day ? _self.day : day // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,total: null == total ? _self.total : total // ignore: cast_nullable_to_non_nullable
as int,sent: null == sent ? _self.sent : sent // ignore: cast_nullable_to_non_nullable
as int,delivered: null == delivered ? _self.delivered : delivered // ignore: cast_nullable_to_non_nullable
as int,failed: null == failed ? _self.failed : failed // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [DailyPoint].
extension DailyPointPatterns on DailyPoint {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _DailyPoint value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _DailyPoint() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _DailyPoint value)  $default,){
final _that = this;
switch (_that) {
case _DailyPoint():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _DailyPoint value)?  $default,){
final _that = this;
switch (_that) {
case _DailyPoint() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String day,  String channel,  int total,  int sent,  int delivered,  int failed)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _DailyPoint() when $default != null:
return $default(_that.day,_that.channel,_that.total,_that.sent,_that.delivered,_that.failed);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String day,  String channel,  int total,  int sent,  int delivered,  int failed)  $default,) {final _that = this;
switch (_that) {
case _DailyPoint():
return $default(_that.day,_that.channel,_that.total,_that.sent,_that.delivered,_that.failed);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String day,  String channel,  int total,  int sent,  int delivered,  int failed)?  $default,) {final _that = this;
switch (_that) {
case _DailyPoint() when $default != null:
return $default(_that.day,_that.channel,_that.total,_that.sent,_that.delivered,_that.failed);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _DailyPoint implements DailyPoint {
  const _DailyPoint({required this.day, required this.channel, this.total = 0, this.sent = 0, this.delivered = 0, this.failed = 0});
  factory _DailyPoint.fromJson(Map<String, dynamic> json) => _$DailyPointFromJson(json);

@override final  String day;
@override final  String channel;
@override@JsonKey() final  int total;
@override@JsonKey() final  int sent;
@override@JsonKey() final  int delivered;
@override@JsonKey() final  int failed;

/// Create a copy of DailyPoint
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$DailyPointCopyWith<_DailyPoint> get copyWith => __$DailyPointCopyWithImpl<_DailyPoint>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$DailyPointToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _DailyPoint&&(identical(other.day, day) || other.day == day)&&(identical(other.channel, channel) || other.channel == channel)&&(identical(other.total, total) || other.total == total)&&(identical(other.sent, sent) || other.sent == sent)&&(identical(other.delivered, delivered) || other.delivered == delivered)&&(identical(other.failed, failed) || other.failed == failed));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,day,channel,total,sent,delivered,failed);
}

@override
String toString() {
    return 'DailyPoint(day: $day, channel: $channel, total: $total, sent: $sent, delivered: $delivered, failed: $failed)';
}


}

/// @nodoc
abstract mixin class _$DailyPointCopyWith<$Res> implements $DailyPointCopyWith<$Res> {
  factory _$DailyPointCopyWith(_DailyPoint value, $Res Function(_DailyPoint) _then) = __$DailyPointCopyWithImpl;
@override @useResult
$Res call({
 String day, String channel, int total, int sent, int delivered, int failed
});




}
/// @nodoc
class __$DailyPointCopyWithImpl<$Res>
    implements _$DailyPointCopyWith<$Res> {
  __$DailyPointCopyWithImpl(this._self, this._then);

  final _DailyPoint _self;
  final $Res Function(_DailyPoint) _then;

/// Create a copy of DailyPoint
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? day = null,Object? channel = null,Object? total = null,Object? sent = null,Object? delivered = null,Object? failed = null,}) {
  return _then(_DailyPoint(
day: null == day ? _self.day : day // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,total: null == total ? _self.total : total // ignore: cast_nullable_to_non_nullable
as int,sent: null == sent ? _self.sent : sent // ignore: cast_nullable_to_non_nullable
as int,delivered: null == delivered ? _self.delivered : delivered // ignore: cast_nullable_to_non_nullable
as int,failed: null == failed ? _self.failed : failed // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$Latency {

@JsonKey(name: 'p50_sent_sec') double get p50SentSec;@JsonKey(name: 'p95_sent_sec') double get p95SentSec;@JsonKey(name: 'p95_delivered_sec') double get p95DeliveredSec; int get samples;
/// Create a copy of Latency
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$LatencyCopyWith<Latency> get copyWith => _$LatencyCopyWithImpl<Latency>(this as Latency, _$identity);

  /// Serializes this Latency to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Latency;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Latency&&(identical(other.p50SentSec, _this.p50SentSec) || other.p50SentSec == _this.p50SentSec)&&(identical(other.p95SentSec, _this.p95SentSec) || other.p95SentSec == _this.p95SentSec)&&(identical(other.p95DeliveredSec, _this.p95DeliveredSec) || other.p95DeliveredSec == _this.p95DeliveredSec)&&(identical(other.samples, _this.samples) || other.samples == _this.samples));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Latency;
  return Object.hash(runtimeType,_this.p50SentSec,_this.p95SentSec,_this.p95DeliveredSec,_this.samples);
}

@override
String toString() {
  final _this = this as Latency;
  return 'Latency(p50SentSec: ${_this.p50SentSec}, p95SentSec: ${_this.p95SentSec}, p95DeliveredSec: ${_this.p95DeliveredSec}, samples: ${_this.samples})';
}


}

/// @nodoc
abstract mixin class $LatencyCopyWith<$Res>  {
  factory $LatencyCopyWith(Latency value, $Res Function(Latency) _then) = _$LatencyCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'p50_sent_sec') double p50SentSec,@JsonKey(name: 'p95_sent_sec') double p95SentSec,@JsonKey(name: 'p95_delivered_sec') double p95DeliveredSec, int samples
});




}
/// @nodoc
class _$LatencyCopyWithImpl<$Res>
    implements $LatencyCopyWith<$Res> {
  _$LatencyCopyWithImpl(this._self, this._then);

  final Latency _self;
  final $Res Function(Latency) _then;

/// Create a copy of Latency
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? p50SentSec = null,Object? p95SentSec = null,Object? p95DeliveredSec = null,Object? samples = null,}) {
  return _then(Latency(
p50SentSec: null == p50SentSec ? _self.p50SentSec : p50SentSec // ignore: cast_nullable_to_non_nullable
as double,p95SentSec: null == p95SentSec ? _self.p95SentSec : p95SentSec // ignore: cast_nullable_to_non_nullable
as double,p95DeliveredSec: null == p95DeliveredSec ? _self.p95DeliveredSec : p95DeliveredSec // ignore: cast_nullable_to_non_nullable
as double,samples: null == samples ? _self.samples : samples // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [Latency].
extension LatencyPatterns on Latency {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Latency value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Latency() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Latency value)  $default,){
final _that = this;
switch (_that) {
case _Latency():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Latency value)?  $default,){
final _that = this;
switch (_that) {
case _Latency() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'p50_sent_sec')  double p50SentSec, @JsonKey(name: 'p95_sent_sec')  double p95SentSec, @JsonKey(name: 'p95_delivered_sec')  double p95DeliveredSec,  int samples)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Latency() when $default != null:
return $default(_that.p50SentSec,_that.p95SentSec,_that.p95DeliveredSec,_that.samples);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'p50_sent_sec')  double p50SentSec, @JsonKey(name: 'p95_sent_sec')  double p95SentSec, @JsonKey(name: 'p95_delivered_sec')  double p95DeliveredSec,  int samples)  $default,) {final _that = this;
switch (_that) {
case _Latency():
return $default(_that.p50SentSec,_that.p95SentSec,_that.p95DeliveredSec,_that.samples);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'p50_sent_sec')  double p50SentSec, @JsonKey(name: 'p95_sent_sec')  double p95SentSec, @JsonKey(name: 'p95_delivered_sec')  double p95DeliveredSec,  int samples)?  $default,) {final _that = this;
switch (_that) {
case _Latency() when $default != null:
return $default(_that.p50SentSec,_that.p95SentSec,_that.p95DeliveredSec,_that.samples);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Latency implements Latency {
  const _Latency({@JsonKey(name: 'p50_sent_sec') this.p50SentSec = 0, @JsonKey(name: 'p95_sent_sec') this.p95SentSec = 0, @JsonKey(name: 'p95_delivered_sec') this.p95DeliveredSec = 0, this.samples = 0});
  factory _Latency.fromJson(Map<String, dynamic> json) => _$LatencyFromJson(json);

@override@JsonKey(name: 'p50_sent_sec') final  double p50SentSec;
@override@JsonKey(name: 'p95_sent_sec') final  double p95SentSec;
@override@JsonKey(name: 'p95_delivered_sec') final  double p95DeliveredSec;
@override@JsonKey() final  int samples;

/// Create a copy of Latency
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$LatencyCopyWith<_Latency> get copyWith => __$LatencyCopyWithImpl<_Latency>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$LatencyToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Latency&&(identical(other.p50SentSec, p50SentSec) || other.p50SentSec == p50SentSec)&&(identical(other.p95SentSec, p95SentSec) || other.p95SentSec == p95SentSec)&&(identical(other.p95DeliveredSec, p95DeliveredSec) || other.p95DeliveredSec == p95DeliveredSec)&&(identical(other.samples, samples) || other.samples == samples));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,p50SentSec,p95SentSec,p95DeliveredSec,samples);
}

@override
String toString() {
    return 'Latency(p50SentSec: $p50SentSec, p95SentSec: $p95SentSec, p95DeliveredSec: $p95DeliveredSec, samples: $samples)';
}


}

/// @nodoc
abstract mixin class _$LatencyCopyWith<$Res> implements $LatencyCopyWith<$Res> {
  factory _$LatencyCopyWith(_Latency value, $Res Function(_Latency) _then) = __$LatencyCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'p50_sent_sec') double p50SentSec,@JsonKey(name: 'p95_sent_sec') double p95SentSec,@JsonKey(name: 'p95_delivered_sec') double p95DeliveredSec, int samples
});




}
/// @nodoc
class __$LatencyCopyWithImpl<$Res>
    implements _$LatencyCopyWith<$Res> {
  __$LatencyCopyWithImpl(this._self, this._then);

  final _Latency _self;
  final $Res Function(_Latency) _then;

/// Create a copy of Latency
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? p50SentSec = null,Object? p95SentSec = null,Object? p95DeliveredSec = null,Object? samples = null,}) {
  return _then(_Latency(
p50SentSec: null == p50SentSec ? _self.p50SentSec : p50SentSec // ignore: cast_nullable_to_non_nullable
as double,p95SentSec: null == p95SentSec ? _self.p95SentSec : p95SentSec // ignore: cast_nullable_to_non_nullable
as double,p95DeliveredSec: null == p95DeliveredSec ? _self.p95DeliveredSec : p95DeliveredSec // ignore: cast_nullable_to_non_nullable
as double,samples: null == samples ? _self.samples : samples // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$Dashboard {

 Totals get totals; Latency get latency; List<DailyPoint> get daily;@JsonKey(name: 'recent_failures') List<Map<String, dynamic>> get recentFailures;@JsonKey(name: 'by_channel') Map<String, Totals> get byChannel;
/// Create a copy of Dashboard
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$DashboardCopyWith<Dashboard> get copyWith => _$DashboardCopyWithImpl<Dashboard>(this as Dashboard, _$identity);

  /// Serializes this Dashboard to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Dashboard;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Dashboard&&(identical(other.totals, _this.totals) || other.totals == _this.totals)&&(identical(other.latency, _this.latency) || other.latency == _this.latency)&&const DeepCollectionEquality().equals(other.daily, _this.daily)&&const DeepCollectionEquality().equals(other.recentFailures, _this.recentFailures)&&const DeepCollectionEquality().equals(other.byChannel, _this.byChannel));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Dashboard;
  return Object.hash(runtimeType,_this.totals,_this.latency,const DeepCollectionEquality().hash(_this.daily),const DeepCollectionEquality().hash(_this.recentFailures),const DeepCollectionEquality().hash(_this.byChannel));
}

@override
String toString() {
  final _this = this as Dashboard;
  return 'Dashboard(totals: ${_this.totals}, latency: ${_this.latency}, daily: ${_this.daily}, recentFailures: ${_this.recentFailures}, byChannel: ${_this.byChannel})';
}


}

/// @nodoc
abstract mixin class $DashboardCopyWith<$Res>  {
  factory $DashboardCopyWith(Dashboard value, $Res Function(Dashboard) _then) = _$DashboardCopyWithImpl;
@useResult
$Res call({
 Totals totals, Latency latency, List<DailyPoint> daily,@JsonKey(name: 'recent_failures') List<Map<String, dynamic>> recentFailures,@JsonKey(name: 'by_channel') Map<String, Totals> byChannel
});


$TotalsCopyWith<$Res> get totals;$LatencyCopyWith<$Res> get latency;

}
/// @nodoc
class _$DashboardCopyWithImpl<$Res>
    implements $DashboardCopyWith<$Res> {
  _$DashboardCopyWithImpl(this._self, this._then);

  final Dashboard _self;
  final $Res Function(Dashboard) _then;

/// Create a copy of Dashboard
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? totals = null,Object? latency = null,Object? daily = null,Object? recentFailures = null,Object? byChannel = null,}) {
  return _then(Dashboard(
totals: null == totals ? _self.totals : totals // ignore: cast_nullable_to_non_nullable
as Totals,latency: null == latency ? _self.latency : latency // ignore: cast_nullable_to_non_nullable
as Latency,daily: null == daily ? _self.daily : daily // ignore: cast_nullable_to_non_nullable
as List<DailyPoint>,recentFailures: null == recentFailures ? _self.recentFailures : recentFailures // ignore: cast_nullable_to_non_nullable
as List<Map<String, dynamic>>,byChannel: null == byChannel ? _self.byChannel : byChannel // ignore: cast_nullable_to_non_nullable
as Map<String, Totals>,
  ));
}
/// Create a copy of Dashboard
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TotalsCopyWith<$Res> get totals {
  
  return $TotalsCopyWith<$Res>(_self.totals, (value) {
    return _then(_self.copyWith(totals: value));
  });
}/// Create a copy of Dashboard
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$LatencyCopyWith<$Res> get latency {
  
  return $LatencyCopyWith<$Res>(_self.latency, (value) {
    return _then(_self.copyWith(latency: value));
  });
}
}


/// Adds pattern-matching-related methods to [Dashboard].
extension DashboardPatterns on Dashboard {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Dashboard value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Dashboard() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Dashboard value)  $default,){
final _that = this;
switch (_that) {
case _Dashboard():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Dashboard value)?  $default,){
final _that = this;
switch (_that) {
case _Dashboard() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( Totals totals,  Latency latency,  List<DailyPoint> daily, @JsonKey(name: 'recent_failures')  List<Map<String, dynamic>> recentFailures, @JsonKey(name: 'by_channel')  Map<String, Totals> byChannel)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Dashboard() when $default != null:
return $default(_that.totals,_that.latency,_that.daily,_that.recentFailures,_that.byChannel);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( Totals totals,  Latency latency,  List<DailyPoint> daily, @JsonKey(name: 'recent_failures')  List<Map<String, dynamic>> recentFailures, @JsonKey(name: 'by_channel')  Map<String, Totals> byChannel)  $default,) {final _that = this;
switch (_that) {
case _Dashboard():
return $default(_that.totals,_that.latency,_that.daily,_that.recentFailures,_that.byChannel);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( Totals totals,  Latency latency,  List<DailyPoint> daily, @JsonKey(name: 'recent_failures')  List<Map<String, dynamic>> recentFailures, @JsonKey(name: 'by_channel')  Map<String, Totals> byChannel)?  $default,) {final _that = this;
switch (_that) {
case _Dashboard() when $default != null:
return $default(_that.totals,_that.latency,_that.daily,_that.recentFailures,_that.byChannel);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Dashboard implements Dashboard {
  const _Dashboard({this.totals = const Totals(), this.latency = const Latency(),  List<DailyPoint> daily = const <DailyPoint>[], @JsonKey(name: 'recent_failures')  List<Map<String, dynamic>> recentFailures = const <Map<String, dynamic>>[], @JsonKey(name: 'by_channel')  Map<String, Totals> byChannel = const <String, Totals>{}}): _daily = daily,_recentFailures = recentFailures,_byChannel = byChannel;
  factory _Dashboard.fromJson(Map<String, dynamic> json) => _$DashboardFromJson(json);

@override@JsonKey() final  Totals totals;
@override@JsonKey() final  Latency latency;
 final  List<DailyPoint> _daily;
@override@JsonKey() List<DailyPoint> get daily {
  if (_daily is EqualUnmodifiableListView) return _daily;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_daily);
}

 final  List<Map<String, dynamic>> _recentFailures;
@override@JsonKey(name: 'recent_failures') List<Map<String, dynamic>> get recentFailures {
  if (_recentFailures is EqualUnmodifiableListView) return _recentFailures;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_recentFailures);
}

 final  Map<String, Totals> _byChannel;
@override@JsonKey(name: 'by_channel') Map<String, Totals> get byChannel {
  if (_byChannel is EqualUnmodifiableMapView) return _byChannel;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableMapView(_byChannel);
}


/// Create a copy of Dashboard
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$DashboardCopyWith<_Dashboard> get copyWith => __$DashboardCopyWithImpl<_Dashboard>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$DashboardToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Dashboard&&(identical(other.totals, totals) || other.totals == totals)&&(identical(other.latency, latency) || other.latency == latency)&&const DeepCollectionEquality().equals(other.daily, _daily)&&const DeepCollectionEquality().equals(other.recentFailures, _recentFailures)&&const DeepCollectionEquality().equals(other.byChannel, _byChannel));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,totals,latency,const DeepCollectionEquality().hash(_daily),const DeepCollectionEquality().hash(_recentFailures),const DeepCollectionEquality().hash(_byChannel));
}

@override
String toString() {
    return 'Dashboard(totals: $totals, latency: $latency, daily: $daily, recentFailures: $recentFailures, byChannel: $byChannel)';
}


}

/// @nodoc
abstract mixin class _$DashboardCopyWith<$Res> implements $DashboardCopyWith<$Res> {
  factory _$DashboardCopyWith(_Dashboard value, $Res Function(_Dashboard) _then) = __$DashboardCopyWithImpl;
@override @useResult
$Res call({
 Totals totals, Latency latency, List<DailyPoint> daily,@JsonKey(name: 'recent_failures') List<Map<String, dynamic>> recentFailures,@JsonKey(name: 'by_channel') Map<String, Totals> byChannel
});


@override $TotalsCopyWith<$Res> get totals;@override $LatencyCopyWith<$Res> get latency;

}
/// @nodoc
class __$DashboardCopyWithImpl<$Res>
    implements _$DashboardCopyWith<$Res> {
  __$DashboardCopyWithImpl(this._self, this._then);

  final _Dashboard _self;
  final $Res Function(_Dashboard) _then;

/// Create a copy of Dashboard
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? totals = null,Object? latency = null,Object? daily = null,Object? recentFailures = null,Object? byChannel = null,}) {
  return _then(_Dashboard(
totals: null == totals ? _self.totals : totals // ignore: cast_nullable_to_non_nullable
as Totals,latency: null == latency ? _self.latency : latency // ignore: cast_nullable_to_non_nullable
as Latency,daily: null == daily ? _self._daily : daily // ignore: cast_nullable_to_non_nullable
as List<DailyPoint>,recentFailures: null == recentFailures ? _self._recentFailures : recentFailures // ignore: cast_nullable_to_non_nullable
as List<Map<String, dynamic>>,byChannel: null == byChannel ? _self._byChannel : byChannel // ignore: cast_nullable_to_non_nullable
as Map<String, Totals>,
  ));
}

/// Create a copy of Dashboard
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TotalsCopyWith<$Res> get totals {
  
  return $TotalsCopyWith<$Res>(_self.totals, (value) {
    return _then(_self.copyWith(totals: value));
  });
}/// Create a copy of Dashboard
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$LatencyCopyWith<$Res> get latency {
  
  return $LatencyCopyWith<$Res>(_self.latency, (value) {
    return _then(_self.copyWith(latency: value));
  });
}
}


/// @nodoc
mixin _$Device {

 String get id; String get platform;@JsonKey(name: 'token_hint') String get tokenHint;@JsonKey(name: 'is_active') bool get isActive;
/// Create a copy of Device
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$DeviceCopyWith<Device> get copyWith => _$DeviceCopyWithImpl<Device>(this as Device, _$identity);

  /// Serializes this Device to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Device;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Device&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.platform, _this.platform) || other.platform == _this.platform)&&(identical(other.tokenHint, _this.tokenHint) || other.tokenHint == _this.tokenHint)&&(identical(other.isActive, _this.isActive) || other.isActive == _this.isActive));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Device;
  return Object.hash(runtimeType,_this.id,_this.platform,_this.tokenHint,_this.isActive);
}

@override
String toString() {
  final _this = this as Device;
  return 'Device(id: ${_this.id}, platform: ${_this.platform}, tokenHint: ${_this.tokenHint}, isActive: ${_this.isActive})';
}


}

/// @nodoc
abstract mixin class $DeviceCopyWith<$Res>  {
  factory $DeviceCopyWith(Device value, $Res Function(Device) _then) = _$DeviceCopyWithImpl;
@useResult
$Res call({
 String id, String platform,@JsonKey(name: 'token_hint') String tokenHint,@JsonKey(name: 'is_active') bool isActive
});




}
/// @nodoc
class _$DeviceCopyWithImpl<$Res>
    implements $DeviceCopyWith<$Res> {
  _$DeviceCopyWithImpl(this._self, this._then);

  final Device _self;
  final $Res Function(Device) _then;

/// Create a copy of Device
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? platform = null,Object? tokenHint = null,Object? isActive = null,}) {
  return _then(Device(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,platform: null == platform ? _self.platform : platform // ignore: cast_nullable_to_non_nullable
as String,tokenHint: null == tokenHint ? _self.tokenHint : tokenHint // ignore: cast_nullable_to_non_nullable
as String,isActive: null == isActive ? _self.isActive : isActive // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

}


/// Adds pattern-matching-related methods to [Device].
extension DevicePatterns on Device {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Device value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Device() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Device value)  $default,){
final _that = this;
switch (_that) {
case _Device():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Device value)?  $default,){
final _that = this;
switch (_that) {
case _Device() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String platform, @JsonKey(name: 'token_hint')  String tokenHint, @JsonKey(name: 'is_active')  bool isActive)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Device() when $default != null:
return $default(_that.id,_that.platform,_that.tokenHint,_that.isActive);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String platform, @JsonKey(name: 'token_hint')  String tokenHint, @JsonKey(name: 'is_active')  bool isActive)  $default,) {final _that = this;
switch (_that) {
case _Device():
return $default(_that.id,_that.platform,_that.tokenHint,_that.isActive);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String platform, @JsonKey(name: 'token_hint')  String tokenHint, @JsonKey(name: 'is_active')  bool isActive)?  $default,) {final _that = this;
switch (_that) {
case _Device() when $default != null:
return $default(_that.id,_that.platform,_that.tokenHint,_that.isActive);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Device implements Device {
  const _Device({required this.id, required this.platform, @JsonKey(name: 'token_hint') this.tokenHint = '', @JsonKey(name: 'is_active') this.isActive = true});
  factory _Device.fromJson(Map<String, dynamic> json) => _$DeviceFromJson(json);

@override final  String id;
@override final  String platform;
@override@JsonKey(name: 'token_hint') final  String tokenHint;
@override@JsonKey(name: 'is_active') final  bool isActive;

/// Create a copy of Device
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$DeviceCopyWith<_Device> get copyWith => __$DeviceCopyWithImpl<_Device>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$DeviceToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Device&&(identical(other.id, id) || other.id == id)&&(identical(other.platform, platform) || other.platform == platform)&&(identical(other.tokenHint, tokenHint) || other.tokenHint == tokenHint)&&(identical(other.isActive, isActive) || other.isActive == isActive));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,platform,tokenHint,isActive);
}

@override
String toString() {
    return 'Device(id: $id, platform: $platform, tokenHint: $tokenHint, isActive: $isActive)';
}


}

/// @nodoc
abstract mixin class _$DeviceCopyWith<$Res> implements $DeviceCopyWith<$Res> {
  factory _$DeviceCopyWith(_Device value, $Res Function(_Device) _then) = __$DeviceCopyWithImpl;
@override @useResult
$Res call({
 String id, String platform,@JsonKey(name: 'token_hint') String tokenHint,@JsonKey(name: 'is_active') bool isActive
});




}
/// @nodoc
class __$DeviceCopyWithImpl<$Res>
    implements _$DeviceCopyWith<$Res> {
  __$DeviceCopyWithImpl(this._self, this._then);

  final _Device _self;
  final $Res Function(_Device) _then;

/// Create a copy of Device
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? platform = null,Object? tokenHint = null,Object? isActive = null,}) {
  return _then(_Device(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,platform: null == platform ? _self.platform : platform // ignore: cast_nullable_to_non_nullable
as String,tokenHint: null == tokenHint ? _self.tokenHint : tokenHint // ignore: cast_nullable_to_non_nullable
as String,isActive: null == isActive ? _self.isActive : isActive // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}


}


/// @nodoc
mixin _$Contact {

 String get id;@JsonKey(name: 'external_id') String get externalId; String get name; String get phone; String get email; String get whatsapp;@JsonKey(name: 'telegram_chat_id') String get telegramChatId;@JsonKey(name: 'slack_id') String get slackId; String get locale; List<String> get tags;@JsonKey(name: 'created_at') DateTime? get createdAt;
/// Create a copy of Contact
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ContactCopyWith<Contact> get copyWith => _$ContactCopyWithImpl<Contact>(this as Contact, _$identity);

  /// Serializes this Contact to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Contact;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Contact&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.externalId, _this.externalId) || other.externalId == _this.externalId)&&(identical(other.name, _this.name) || other.name == _this.name)&&(identical(other.phone, _this.phone) || other.phone == _this.phone)&&(identical(other.email, _this.email) || other.email == _this.email)&&(identical(other.whatsapp, _this.whatsapp) || other.whatsapp == _this.whatsapp)&&(identical(other.telegramChatId, _this.telegramChatId) || other.telegramChatId == _this.telegramChatId)&&(identical(other.slackId, _this.slackId) || other.slackId == _this.slackId)&&(identical(other.locale, _this.locale) || other.locale == _this.locale)&&const DeepCollectionEquality().equals(other.tags, _this.tags)&&(identical(other.createdAt, _this.createdAt) || other.createdAt == _this.createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Contact;
  return Object.hash(runtimeType,_this.id,_this.externalId,_this.name,_this.phone,_this.email,_this.whatsapp,_this.telegramChatId,_this.slackId,_this.locale,const DeepCollectionEquality().hash(_this.tags),_this.createdAt);
}

@override
String toString() {
  final _this = this as Contact;
  return 'Contact(id: ${_this.id}, externalId: ${_this.externalId}, name: ${_this.name}, phone: ${_this.phone}, email: ${_this.email}, whatsapp: ${_this.whatsapp}, telegramChatId: ${_this.telegramChatId}, slackId: ${_this.slackId}, locale: ${_this.locale}, tags: ${_this.tags}, createdAt: ${_this.createdAt})';
}


}

/// @nodoc
abstract mixin class $ContactCopyWith<$Res>  {
  factory $ContactCopyWith(Contact value, $Res Function(Contact) _then) = _$ContactCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'external_id') String externalId, String name, String phone, String email, String whatsapp,@JsonKey(name: 'telegram_chat_id') String telegramChatId,@JsonKey(name: 'slack_id') String slackId, String locale, List<String> tags,@JsonKey(name: 'created_at') DateTime? createdAt
});




}
/// @nodoc
class _$ContactCopyWithImpl<$Res>
    implements $ContactCopyWith<$Res> {
  _$ContactCopyWithImpl(this._self, this._then);

  final Contact _self;
  final $Res Function(Contact) _then;

/// Create a copy of Contact
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? externalId = null,Object? name = null,Object? phone = null,Object? email = null,Object? whatsapp = null,Object? telegramChatId = null,Object? slackId = null,Object? locale = null,Object? tags = null,Object? createdAt = freezed,}) {
  return _then(Contact(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,externalId: null == externalId ? _self.externalId : externalId // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,phone: null == phone ? _self.phone : phone // ignore: cast_nullable_to_non_nullable
as String,email: null == email ? _self.email : email // ignore: cast_nullable_to_non_nullable
as String,whatsapp: null == whatsapp ? _self.whatsapp : whatsapp // ignore: cast_nullable_to_non_nullable
as String,telegramChatId: null == telegramChatId ? _self.telegramChatId : telegramChatId // ignore: cast_nullable_to_non_nullable
as String,slackId: null == slackId ? _self.slackId : slackId // ignore: cast_nullable_to_non_nullable
as String,locale: null == locale ? _self.locale : locale // ignore: cast_nullable_to_non_nullable
as String,tags: null == tags ? _self.tags : tags // ignore: cast_nullable_to_non_nullable
as List<String>,createdAt: freezed == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime?,
  ));
}

}


/// Adds pattern-matching-related methods to [Contact].
extension ContactPatterns on Contact {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Contact value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Contact() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Contact value)  $default,){
final _that = this;
switch (_that) {
case _Contact():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Contact value)?  $default,){
final _that = this;
switch (_that) {
case _Contact() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'external_id')  String externalId,  String name,  String phone,  String email,  String whatsapp, @JsonKey(name: 'telegram_chat_id')  String telegramChatId, @JsonKey(name: 'slack_id')  String slackId,  String locale,  List<String> tags, @JsonKey(name: 'created_at')  DateTime? createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Contact() when $default != null:
return $default(_that.id,_that.externalId,_that.name,_that.phone,_that.email,_that.whatsapp,_that.telegramChatId,_that.slackId,_that.locale,_that.tags,_that.createdAt);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'external_id')  String externalId,  String name,  String phone,  String email,  String whatsapp, @JsonKey(name: 'telegram_chat_id')  String telegramChatId, @JsonKey(name: 'slack_id')  String slackId,  String locale,  List<String> tags, @JsonKey(name: 'created_at')  DateTime? createdAt)  $default,) {final _that = this;
switch (_that) {
case _Contact():
return $default(_that.id,_that.externalId,_that.name,_that.phone,_that.email,_that.whatsapp,_that.telegramChatId,_that.slackId,_that.locale,_that.tags,_that.createdAt);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'external_id')  String externalId,  String name,  String phone,  String email,  String whatsapp, @JsonKey(name: 'telegram_chat_id')  String telegramChatId, @JsonKey(name: 'slack_id')  String slackId,  String locale,  List<String> tags, @JsonKey(name: 'created_at')  DateTime? createdAt)?  $default,) {final _that = this;
switch (_that) {
case _Contact() when $default != null:
return $default(_that.id,_that.externalId,_that.name,_that.phone,_that.email,_that.whatsapp,_that.telegramChatId,_that.slackId,_that.locale,_that.tags,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Contact implements Contact {
  const _Contact({required this.id, @JsonKey(name: 'external_id') this.externalId = '', this.name = '', this.phone = '', this.email = '', this.whatsapp = '', @JsonKey(name: 'telegram_chat_id') this.telegramChatId = '', @JsonKey(name: 'slack_id') this.slackId = '', this.locale = 'tk',  List<String> tags = const <String>[], @JsonKey(name: 'created_at') this.createdAt}): _tags = tags;
  factory _Contact.fromJson(Map<String, dynamic> json) => _$ContactFromJson(json);

@override final  String id;
@override@JsonKey(name: 'external_id') final  String externalId;
@override@JsonKey() final  String name;
@override@JsonKey() final  String phone;
@override@JsonKey() final  String email;
@override@JsonKey() final  String whatsapp;
@override@JsonKey(name: 'telegram_chat_id') final  String telegramChatId;
@override@JsonKey(name: 'slack_id') final  String slackId;
@override@JsonKey() final  String locale;
 final  List<String> _tags;
@override@JsonKey() List<String> get tags {
  if (_tags is EqualUnmodifiableListView) return _tags;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_tags);
}

@override@JsonKey(name: 'created_at') final  DateTime? createdAt;

/// Create a copy of Contact
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ContactCopyWith<_Contact> get copyWith => __$ContactCopyWithImpl<_Contact>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ContactToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Contact&&(identical(other.id, id) || other.id == id)&&(identical(other.externalId, externalId) || other.externalId == externalId)&&(identical(other.name, name) || other.name == name)&&(identical(other.phone, phone) || other.phone == phone)&&(identical(other.email, email) || other.email == email)&&(identical(other.whatsapp, whatsapp) || other.whatsapp == whatsapp)&&(identical(other.telegramChatId, telegramChatId) || other.telegramChatId == telegramChatId)&&(identical(other.slackId, slackId) || other.slackId == slackId)&&(identical(other.locale, locale) || other.locale == locale)&&const DeepCollectionEquality().equals(other.tags, _tags)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,externalId,name,phone,email,whatsapp,telegramChatId,slackId,locale,const DeepCollectionEquality().hash(_tags),createdAt);
}

@override
String toString() {
    return 'Contact(id: $id, externalId: $externalId, name: $name, phone: $phone, email: $email, whatsapp: $whatsapp, telegramChatId: $telegramChatId, slackId: $slackId, locale: $locale, tags: $tags, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$ContactCopyWith<$Res> implements $ContactCopyWith<$Res> {
  factory _$ContactCopyWith(_Contact value, $Res Function(_Contact) _then) = __$ContactCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'external_id') String externalId, String name, String phone, String email, String whatsapp,@JsonKey(name: 'telegram_chat_id') String telegramChatId,@JsonKey(name: 'slack_id') String slackId, String locale, List<String> tags,@JsonKey(name: 'created_at') DateTime? createdAt
});




}
/// @nodoc
class __$ContactCopyWithImpl<$Res>
    implements _$ContactCopyWith<$Res> {
  __$ContactCopyWithImpl(this._self, this._then);

  final _Contact _self;
  final $Res Function(_Contact) _then;

/// Create a copy of Contact
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? externalId = null,Object? name = null,Object? phone = null,Object? email = null,Object? whatsapp = null,Object? telegramChatId = null,Object? slackId = null,Object? locale = null,Object? tags = null,Object? createdAt = freezed,}) {
  return _then(_Contact(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,externalId: null == externalId ? _self.externalId : externalId // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,phone: null == phone ? _self.phone : phone // ignore: cast_nullable_to_non_nullable
as String,email: null == email ? _self.email : email // ignore: cast_nullable_to_non_nullable
as String,whatsapp: null == whatsapp ? _self.whatsapp : whatsapp // ignore: cast_nullable_to_non_nullable
as String,telegramChatId: null == telegramChatId ? _self.telegramChatId : telegramChatId // ignore: cast_nullable_to_non_nullable
as String,slackId: null == slackId ? _self.slackId : slackId // ignore: cast_nullable_to_non_nullable
as String,locale: null == locale ? _self.locale : locale // ignore: cast_nullable_to_non_nullable
as String,tags: null == tags ? _self._tags : tags // ignore: cast_nullable_to_non_nullable
as List<String>,createdAt: freezed == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime?,
  ));
}


}


/// @nodoc
mixin _$Group {

 String get id; String get name; String get description;@JsonKey(name: 'member_count') int get memberCount;
/// Create a copy of Group
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$GroupCopyWith<Group> get copyWith => _$GroupCopyWithImpl<Group>(this as Group, _$identity);

  /// Serializes this Group to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Group;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Group&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.name, _this.name) || other.name == _this.name)&&(identical(other.description, _this.description) || other.description == _this.description)&&(identical(other.memberCount, _this.memberCount) || other.memberCount == _this.memberCount));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Group;
  return Object.hash(runtimeType,_this.id,_this.name,_this.description,_this.memberCount);
}

@override
String toString() {
  final _this = this as Group;
  return 'Group(id: ${_this.id}, name: ${_this.name}, description: ${_this.description}, memberCount: ${_this.memberCount})';
}


}

/// @nodoc
abstract mixin class $GroupCopyWith<$Res>  {
  factory $GroupCopyWith(Group value, $Res Function(Group) _then) = _$GroupCopyWithImpl;
@useResult
$Res call({
 String id, String name, String description,@JsonKey(name: 'member_count') int memberCount
});




}
/// @nodoc
class _$GroupCopyWithImpl<$Res>
    implements $GroupCopyWith<$Res> {
  _$GroupCopyWithImpl(this._self, this._then);

  final Group _self;
  final $Res Function(Group) _then;

/// Create a copy of Group
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? name = null,Object? description = null,Object? memberCount = null,}) {
  return _then(Group(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,description: null == description ? _self.description : description // ignore: cast_nullable_to_non_nullable
as String,memberCount: null == memberCount ? _self.memberCount : memberCount // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [Group].
extension GroupPatterns on Group {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Group value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Group() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Group value)  $default,){
final _that = this;
switch (_that) {
case _Group():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Group value)?  $default,){
final _that = this;
switch (_that) {
case _Group() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String name,  String description, @JsonKey(name: 'member_count')  int memberCount)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Group() when $default != null:
return $default(_that.id,_that.name,_that.description,_that.memberCount);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String name,  String description, @JsonKey(name: 'member_count')  int memberCount)  $default,) {final _that = this;
switch (_that) {
case _Group():
return $default(_that.id,_that.name,_that.description,_that.memberCount);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String name,  String description, @JsonKey(name: 'member_count')  int memberCount)?  $default,) {final _that = this;
switch (_that) {
case _Group() when $default != null:
return $default(_that.id,_that.name,_that.description,_that.memberCount);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Group implements Group {
  const _Group({required this.id, required this.name, this.description = '', @JsonKey(name: 'member_count') this.memberCount = 0});
  factory _Group.fromJson(Map<String, dynamic> json) => _$GroupFromJson(json);

@override final  String id;
@override final  String name;
@override@JsonKey() final  String description;
@override@JsonKey(name: 'member_count') final  int memberCount;

/// Create a copy of Group
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$GroupCopyWith<_Group> get copyWith => __$GroupCopyWithImpl<_Group>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$GroupToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Group&&(identical(other.id, id) || other.id == id)&&(identical(other.name, name) || other.name == name)&&(identical(other.description, description) || other.description == description)&&(identical(other.memberCount, memberCount) || other.memberCount == memberCount));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,name,description,memberCount);
}

@override
String toString() {
    return 'Group(id: $id, name: $name, description: $description, memberCount: $memberCount)';
}


}

/// @nodoc
abstract mixin class _$GroupCopyWith<$Res> implements $GroupCopyWith<$Res> {
  factory _$GroupCopyWith(_Group value, $Res Function(_Group) _then) = __$GroupCopyWithImpl;
@override @useResult
$Res call({
 String id, String name, String description,@JsonKey(name: 'member_count') int memberCount
});




}
/// @nodoc
class __$GroupCopyWithImpl<$Res>
    implements _$GroupCopyWith<$Res> {
  __$GroupCopyWithImpl(this._self, this._then);

  final _Group _self;
  final $Res Function(_Group) _then;

/// Create a copy of Group
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? name = null,Object? description = null,Object? memberCount = null,}) {
  return _then(_Group(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,description: null == description ? _self.description : description // ignore: cast_nullable_to_non_nullable
as String,memberCount: null == memberCount ? _self.memberCount : memberCount // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$MembersResult {

 int get added;@JsonKey(name: 'created_contacts') int get createdContacts;@JsonKey(name: 'not_found') List<String> get notFound;
/// Create a copy of MembersResult
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MembersResultCopyWith<MembersResult> get copyWith => _$MembersResultCopyWithImpl<MembersResult>(this as MembersResult, _$identity);

  /// Serializes this MembersResult to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as MembersResult;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MembersResult&&(identical(other.added, _this.added) || other.added == _this.added)&&(identical(other.createdContacts, _this.createdContacts) || other.createdContacts == _this.createdContacts)&&const DeepCollectionEquality().equals(other.notFound, _this.notFound));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as MembersResult;
  return Object.hash(runtimeType,_this.added,_this.createdContacts,const DeepCollectionEquality().hash(_this.notFound));
}

@override
String toString() {
  final _this = this as MembersResult;
  return 'MembersResult(added: ${_this.added}, createdContacts: ${_this.createdContacts}, notFound: ${_this.notFound})';
}


}

/// @nodoc
abstract mixin class $MembersResultCopyWith<$Res>  {
  factory $MembersResultCopyWith(MembersResult value, $Res Function(MembersResult) _then) = _$MembersResultCopyWithImpl;
@useResult
$Res call({
 int added,@JsonKey(name: 'created_contacts') int createdContacts,@JsonKey(name: 'not_found') List<String> notFound
});




}
/// @nodoc
class _$MembersResultCopyWithImpl<$Res>
    implements $MembersResultCopyWith<$Res> {
  _$MembersResultCopyWithImpl(this._self, this._then);

  final MembersResult _self;
  final $Res Function(MembersResult) _then;

/// Create a copy of MembersResult
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? added = null,Object? createdContacts = null,Object? notFound = null,}) {
  return _then(MembersResult(
added: null == added ? _self.added : added // ignore: cast_nullable_to_non_nullable
as int,createdContacts: null == createdContacts ? _self.createdContacts : createdContacts // ignore: cast_nullable_to_non_nullable
as int,notFound: null == notFound ? _self.notFound : notFound // ignore: cast_nullable_to_non_nullable
as List<String>,
  ));
}

}


/// Adds pattern-matching-related methods to [MembersResult].
extension MembersResultPatterns on MembersResult {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MembersResult value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MembersResult() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MembersResult value)  $default,){
final _that = this;
switch (_that) {
case _MembersResult():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MembersResult value)?  $default,){
final _that = this;
switch (_that) {
case _MembersResult() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( int added, @JsonKey(name: 'created_contacts')  int createdContacts, @JsonKey(name: 'not_found')  List<String> notFound)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MembersResult() when $default != null:
return $default(_that.added,_that.createdContacts,_that.notFound);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( int added, @JsonKey(name: 'created_contacts')  int createdContacts, @JsonKey(name: 'not_found')  List<String> notFound)  $default,) {final _that = this;
switch (_that) {
case _MembersResult():
return $default(_that.added,_that.createdContacts,_that.notFound);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( int added, @JsonKey(name: 'created_contacts')  int createdContacts, @JsonKey(name: 'not_found')  List<String> notFound)?  $default,) {final _that = this;
switch (_that) {
case _MembersResult() when $default != null:
return $default(_that.added,_that.createdContacts,_that.notFound);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _MembersResult implements MembersResult {
  const _MembersResult({this.added = 0, @JsonKey(name: 'created_contacts') this.createdContacts = 0, @JsonKey(name: 'not_found')  List<String> notFound = const <String>[]}): _notFound = notFound;
  factory _MembersResult.fromJson(Map<String, dynamic> json) => _$MembersResultFromJson(json);

@override@JsonKey() final  int added;
@override@JsonKey(name: 'created_contacts') final  int createdContacts;
 final  List<String> _notFound;
@override@JsonKey(name: 'not_found') List<String> get notFound {
  if (_notFound is EqualUnmodifiableListView) return _notFound;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_notFound);
}


/// Create a copy of MembersResult
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MembersResultCopyWith<_MembersResult> get copyWith => __$MembersResultCopyWithImpl<_MembersResult>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MembersResultToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _MembersResult&&(identical(other.added, added) || other.added == added)&&(identical(other.createdContacts, createdContacts) || other.createdContacts == createdContacts)&&const DeepCollectionEquality().equals(other.notFound, _notFound));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,added,createdContacts,const DeepCollectionEquality().hash(_notFound));
}

@override
String toString() {
    return 'MembersResult(added: $added, createdContacts: $createdContacts, notFound: $notFound)';
}


}

/// @nodoc
abstract mixin class _$MembersResultCopyWith<$Res> implements $MembersResultCopyWith<$Res> {
  factory _$MembersResultCopyWith(_MembersResult value, $Res Function(_MembersResult) _then) = __$MembersResultCopyWithImpl;
@override @useResult
$Res call({
 int added,@JsonKey(name: 'created_contacts') int createdContacts,@JsonKey(name: 'not_found') List<String> notFound
});




}
/// @nodoc
class __$MembersResultCopyWithImpl<$Res>
    implements _$MembersResultCopyWith<$Res> {
  __$MembersResultCopyWithImpl(this._self, this._then);

  final _MembersResult _self;
  final $Res Function(_MembersResult) _then;

/// Create a copy of MembersResult
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? added = null,Object? createdContacts = null,Object? notFound = null,}) {
  return _then(_MembersResult(
added: null == added ? _self.added : added // ignore: cast_nullable_to_non_nullable
as int,createdContacts: null == createdContacts ? _self.createdContacts : createdContacts // ignore: cast_nullable_to_non_nullable
as int,notFound: null == notFound ? _self._notFound : notFound // ignore: cast_nullable_to_non_nullable
as List<String>,
  ));
}


}


/// @nodoc
mixin _$ProviderDetail {

 String get id; String get name; String get channel; String get type; int get priority;@JsonKey(name: 'is_active') bool get isActive;@JsonKey(name: 'rate_limit_per_sec') int get rateLimitPerSec; Map<String, dynamic> get settings;
/// Create a copy of ProviderDetail
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ProviderDetailCopyWith<ProviderDetail> get copyWith => _$ProviderDetailCopyWithImpl<ProviderDetail>(this as ProviderDetail, _$identity);

  /// Serializes this ProviderDetail to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as ProviderDetail;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ProviderDetail&&(identical(other.id, _this.id) || other.id == _this.id)&&(identical(other.name, _this.name) || other.name == _this.name)&&(identical(other.channel, _this.channel) || other.channel == _this.channel)&&(identical(other.type, _this.type) || other.type == _this.type)&&(identical(other.priority, _this.priority) || other.priority == _this.priority)&&(identical(other.isActive, _this.isActive) || other.isActive == _this.isActive)&&(identical(other.rateLimitPerSec, _this.rateLimitPerSec) || other.rateLimitPerSec == _this.rateLimitPerSec)&&const DeepCollectionEquality().equals(other.settings, _this.settings));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as ProviderDetail;
  return Object.hash(runtimeType,_this.id,_this.name,_this.channel,_this.type,_this.priority,_this.isActive,_this.rateLimitPerSec,const DeepCollectionEquality().hash(_this.settings));
}

@override
String toString() {
  final _this = this as ProviderDetail;
  return 'ProviderDetail(id: ${_this.id}, name: ${_this.name}, channel: ${_this.channel}, type: ${_this.type}, priority: ${_this.priority}, isActive: ${_this.isActive}, rateLimitPerSec: ${_this.rateLimitPerSec}, settings: ${_this.settings})';
}


}

/// @nodoc
abstract mixin class $ProviderDetailCopyWith<$Res>  {
  factory $ProviderDetailCopyWith(ProviderDetail value, $Res Function(ProviderDetail) _then) = _$ProviderDetailCopyWithImpl;
@useResult
$Res call({
 String id, String name, String channel, String type, int priority,@JsonKey(name: 'is_active') bool isActive,@JsonKey(name: 'rate_limit_per_sec') int rateLimitPerSec, Map<String, dynamic> settings
});




}
/// @nodoc
class _$ProviderDetailCopyWithImpl<$Res>
    implements $ProviderDetailCopyWith<$Res> {
  _$ProviderDetailCopyWithImpl(this._self, this._then);

  final ProviderDetail _self;
  final $Res Function(ProviderDetail) _then;

/// Create a copy of ProviderDetail
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? name = null,Object? channel = null,Object? type = null,Object? priority = null,Object? isActive = null,Object? rateLimitPerSec = null,Object? settings = null,}) {
  return _then(ProviderDetail(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,priority: null == priority ? _self.priority : priority // ignore: cast_nullable_to_non_nullable
as int,isActive: null == isActive ? _self.isActive : isActive // ignore: cast_nullable_to_non_nullable
as bool,rateLimitPerSec: null == rateLimitPerSec ? _self.rateLimitPerSec : rateLimitPerSec // ignore: cast_nullable_to_non_nullable
as int,settings: null == settings ? _self.settings : settings // ignore: cast_nullable_to_non_nullable
as Map<String, dynamic>,
  ));
}

}


/// Adds pattern-matching-related methods to [ProviderDetail].
extension ProviderDetailPatterns on ProviderDetail {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ProviderDetail value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ProviderDetail() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ProviderDetail value)  $default,){
final _that = this;
switch (_that) {
case _ProviderDetail():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ProviderDetail value)?  $default,){
final _that = this;
switch (_that) {
case _ProviderDetail() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String name,  String channel,  String type,  int priority, @JsonKey(name: 'is_active')  bool isActive, @JsonKey(name: 'rate_limit_per_sec')  int rateLimitPerSec,  Map<String, dynamic> settings)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ProviderDetail() when $default != null:
return $default(_that.id,_that.name,_that.channel,_that.type,_that.priority,_that.isActive,_that.rateLimitPerSec,_that.settings);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String name,  String channel,  String type,  int priority, @JsonKey(name: 'is_active')  bool isActive, @JsonKey(name: 'rate_limit_per_sec')  int rateLimitPerSec,  Map<String, dynamic> settings)  $default,) {final _that = this;
switch (_that) {
case _ProviderDetail():
return $default(_that.id,_that.name,_that.channel,_that.type,_that.priority,_that.isActive,_that.rateLimitPerSec,_that.settings);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String name,  String channel,  String type,  int priority, @JsonKey(name: 'is_active')  bool isActive, @JsonKey(name: 'rate_limit_per_sec')  int rateLimitPerSec,  Map<String, dynamic> settings)?  $default,) {final _that = this;
switch (_that) {
case _ProviderDetail() when $default != null:
return $default(_that.id,_that.name,_that.channel,_that.type,_that.priority,_that.isActive,_that.rateLimitPerSec,_that.settings);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ProviderDetail implements ProviderDetail {
  const _ProviderDetail({required this.id, required this.name, required this.channel, required this.type, this.priority = 100, @JsonKey(name: 'is_active') this.isActive = true, @JsonKey(name: 'rate_limit_per_sec') this.rateLimitPerSec = 0,  Map<String, dynamic> settings = const <String, dynamic>{}}): _settings = settings;
  factory _ProviderDetail.fromJson(Map<String, dynamic> json) => _$ProviderDetailFromJson(json);

@override final  String id;
@override final  String name;
@override final  String channel;
@override final  String type;
@override@JsonKey() final  int priority;
@override@JsonKey(name: 'is_active') final  bool isActive;
@override@JsonKey(name: 'rate_limit_per_sec') final  int rateLimitPerSec;
 final  Map<String, dynamic> _settings;
@override@JsonKey() Map<String, dynamic> get settings {
  if (_settings is EqualUnmodifiableMapView) return _settings;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableMapView(_settings);
}


/// Create a copy of ProviderDetail
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ProviderDetailCopyWith<_ProviderDetail> get copyWith => __$ProviderDetailCopyWithImpl<_ProviderDetail>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ProviderDetailToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _ProviderDetail&&(identical(other.id, id) || other.id == id)&&(identical(other.name, name) || other.name == name)&&(identical(other.channel, channel) || other.channel == channel)&&(identical(other.type, type) || other.type == type)&&(identical(other.priority, priority) || other.priority == priority)&&(identical(other.isActive, isActive) || other.isActive == isActive)&&(identical(other.rateLimitPerSec, rateLimitPerSec) || other.rateLimitPerSec == rateLimitPerSec)&&const DeepCollectionEquality().equals(other.settings, _settings));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,id,name,channel,type,priority,isActive,rateLimitPerSec,const DeepCollectionEquality().hash(_settings));
}

@override
String toString() {
    return 'ProviderDetail(id: $id, name: $name, channel: $channel, type: $type, priority: $priority, isActive: $isActive, rateLimitPerSec: $rateLimitPerSec, settings: $settings)';
}


}

/// @nodoc
abstract mixin class _$ProviderDetailCopyWith<$Res> implements $ProviderDetailCopyWith<$Res> {
  factory _$ProviderDetailCopyWith(_ProviderDetail value, $Res Function(_ProviderDetail) _then) = __$ProviderDetailCopyWithImpl;
@override @useResult
$Res call({
 String id, String name, String channel, String type, int priority,@JsonKey(name: 'is_active') bool isActive,@JsonKey(name: 'rate_limit_per_sec') int rateLimitPerSec, Map<String, dynamic> settings
});




}
/// @nodoc
class __$ProviderDetailCopyWithImpl<$Res>
    implements _$ProviderDetailCopyWith<$Res> {
  __$ProviderDetailCopyWithImpl(this._self, this._then);

  final _ProviderDetail _self;
  final $Res Function(_ProviderDetail) _then;

/// Create a copy of ProviderDetail
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? name = null,Object? channel = null,Object? type = null,Object? priority = null,Object? isActive = null,Object? rateLimitPerSec = null,Object? settings = null,}) {
  return _then(_ProviderDetail(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,channel: null == channel ? _self.channel : channel // ignore: cast_nullable_to_non_nullable
as String,type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,priority: null == priority ? _self.priority : priority // ignore: cast_nullable_to_non_nullable
as int,isActive: null == isActive ? _self.isActive : isActive // ignore: cast_nullable_to_non_nullable
as bool,rateLimitPerSec: null == rateLimitPerSec ? _self.rateLimitPerSec : rateLimitPerSec // ignore: cast_nullable_to_non_nullable
as int,settings: null == settings ? _self._settings : settings // ignore: cast_nullable_to_non_nullable
as Map<String, dynamic>,
  ));
}


}


/// @nodoc
mixin _$Pairing {

@JsonKey(name: 'api_url') String get apiUrl;@JsonKey(name: 'gateway_key') String get gatewayKey; String get qr; bool get online;@JsonKey(name: 'last_seen_at') DateTime? get lastSeenAt;
/// Create a copy of Pairing
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PairingCopyWith<Pairing> get copyWith => _$PairingCopyWithImpl<Pairing>(this as Pairing, _$identity);

  /// Serializes this Pairing to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  final _this = this as Pairing;
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Pairing&&(identical(other.apiUrl, _this.apiUrl) || other.apiUrl == _this.apiUrl)&&(identical(other.gatewayKey, _this.gatewayKey) || other.gatewayKey == _this.gatewayKey)&&(identical(other.qr, _this.qr) || other.qr == _this.qr)&&(identical(other.online, _this.online) || other.online == _this.online)&&(identical(other.lastSeenAt, _this.lastSeenAt) || other.lastSeenAt == _this.lastSeenAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
  final _this = this as Pairing;
  return Object.hash(runtimeType,_this.apiUrl,_this.gatewayKey,_this.qr,_this.online,_this.lastSeenAt);
}

@override
String toString() {
  final _this = this as Pairing;
  return 'Pairing(apiUrl: ${_this.apiUrl}, gatewayKey: ${_this.gatewayKey}, qr: ${_this.qr}, online: ${_this.online}, lastSeenAt: ${_this.lastSeenAt})';
}


}

/// @nodoc
abstract mixin class $PairingCopyWith<$Res>  {
  factory $PairingCopyWith(Pairing value, $Res Function(Pairing) _then) = _$PairingCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'api_url') String apiUrl,@JsonKey(name: 'gateway_key') String gatewayKey, String qr, bool online,@JsonKey(name: 'last_seen_at') DateTime? lastSeenAt
});




}
/// @nodoc
class _$PairingCopyWithImpl<$Res>
    implements $PairingCopyWith<$Res> {
  _$PairingCopyWithImpl(this._self, this._then);

  final Pairing _self;
  final $Res Function(Pairing) _then;

/// Create a copy of Pairing
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? apiUrl = null,Object? gatewayKey = null,Object? qr = null,Object? online = null,Object? lastSeenAt = freezed,}) {
  return _then(Pairing(
apiUrl: null == apiUrl ? _self.apiUrl : apiUrl // ignore: cast_nullable_to_non_nullable
as String,gatewayKey: null == gatewayKey ? _self.gatewayKey : gatewayKey // ignore: cast_nullable_to_non_nullable
as String,qr: null == qr ? _self.qr : qr // ignore: cast_nullable_to_non_nullable
as String,online: null == online ? _self.online : online // ignore: cast_nullable_to_non_nullable
as bool,lastSeenAt: freezed == lastSeenAt ? _self.lastSeenAt : lastSeenAt // ignore: cast_nullable_to_non_nullable
as DateTime?,
  ));
}

}


/// Adds pattern-matching-related methods to [Pairing].
extension PairingPatterns on Pairing {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Pairing value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Pairing() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Pairing value)  $default,){
final _that = this;
switch (_that) {
case _Pairing():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Pairing value)?  $default,){
final _that = this;
switch (_that) {
case _Pairing() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'api_url')  String apiUrl, @JsonKey(name: 'gateway_key')  String gatewayKey,  String qr,  bool online, @JsonKey(name: 'last_seen_at')  DateTime? lastSeenAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Pairing() when $default != null:
return $default(_that.apiUrl,_that.gatewayKey,_that.qr,_that.online,_that.lastSeenAt);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'api_url')  String apiUrl, @JsonKey(name: 'gateway_key')  String gatewayKey,  String qr,  bool online, @JsonKey(name: 'last_seen_at')  DateTime? lastSeenAt)  $default,) {final _that = this;
switch (_that) {
case _Pairing():
return $default(_that.apiUrl,_that.gatewayKey,_that.qr,_that.online,_that.lastSeenAt);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'api_url')  String apiUrl, @JsonKey(name: 'gateway_key')  String gatewayKey,  String qr,  bool online, @JsonKey(name: 'last_seen_at')  DateTime? lastSeenAt)?  $default,) {final _that = this;
switch (_that) {
case _Pairing() when $default != null:
return $default(_that.apiUrl,_that.gatewayKey,_that.qr,_that.online,_that.lastSeenAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Pairing implements Pairing {
  const _Pairing({@JsonKey(name: 'api_url') this.apiUrl = '', @JsonKey(name: 'gateway_key') this.gatewayKey = '', this.qr = '', this.online = false, @JsonKey(name: 'last_seen_at') this.lastSeenAt});
  factory _Pairing.fromJson(Map<String, dynamic> json) => _$PairingFromJson(json);

@override@JsonKey(name: 'api_url') final  String apiUrl;
@override@JsonKey(name: 'gateway_key') final  String gatewayKey;
@override@JsonKey() final  String qr;
@override@JsonKey() final  bool online;
@override@JsonKey(name: 'last_seen_at') final  DateTime? lastSeenAt;

/// Create a copy of Pairing
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PairingCopyWith<_Pairing> get copyWith => __$PairingCopyWithImpl<_Pairing>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$PairingToJson(this, );
}

@override
bool operator ==(Object other) {
    return identical(this, other) || (other.runtimeType == runtimeType&&other is _Pairing&&(identical(other.apiUrl, apiUrl) || other.apiUrl == apiUrl)&&(identical(other.gatewayKey, gatewayKey) || other.gatewayKey == gatewayKey)&&(identical(other.qr, qr) || other.qr == qr)&&(identical(other.online, online) || other.online == online)&&(identical(other.lastSeenAt, lastSeenAt) || other.lastSeenAt == lastSeenAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode {
    return Object.hash(runtimeType,apiUrl,gatewayKey,qr,online,lastSeenAt);
}

@override
String toString() {
    return 'Pairing(apiUrl: $apiUrl, gatewayKey: $gatewayKey, qr: $qr, online: $online, lastSeenAt: $lastSeenAt)';
}


}

/// @nodoc
abstract mixin class _$PairingCopyWith<$Res> implements $PairingCopyWith<$Res> {
  factory _$PairingCopyWith(_Pairing value, $Res Function(_Pairing) _then) = __$PairingCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'api_url') String apiUrl,@JsonKey(name: 'gateway_key') String gatewayKey, String qr, bool online,@JsonKey(name: 'last_seen_at') DateTime? lastSeenAt
});




}
/// @nodoc
class __$PairingCopyWithImpl<$Res>
    implements _$PairingCopyWith<$Res> {
  __$PairingCopyWithImpl(this._self, this._then);

  final _Pairing _self;
  final $Res Function(_Pairing) _then;

/// Create a copy of Pairing
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? apiUrl = null,Object? gatewayKey = null,Object? qr = null,Object? online = null,Object? lastSeenAt = freezed,}) {
  return _then(_Pairing(
apiUrl: null == apiUrl ? _self.apiUrl : apiUrl // ignore: cast_nullable_to_non_nullable
as String,gatewayKey: null == gatewayKey ? _self.gatewayKey : gatewayKey // ignore: cast_nullable_to_non_nullable
as String,qr: null == qr ? _self.qr : qr // ignore: cast_nullable_to_non_nullable
as String,online: null == online ? _self.online : online // ignore: cast_nullable_to_non_nullable
as bool,lastSeenAt: freezed == lastSeenAt ? _self.lastSeenAt : lastSeenAt // ignore: cast_nullable_to_non_nullable
as DateTime?,
  ));
}


}

// dart format on
