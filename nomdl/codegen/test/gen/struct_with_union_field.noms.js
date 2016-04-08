// This file was generated by nomdl/codegen.
// @flow
/* eslint-disable */

import {
  Field as _Field,
  Kind as _Kind,
  Package as _Package,
  blobType as _blobType,
  createStructClass as _createStructClass,
  float32Type as _float32Type,
  float64Type as _float64Type,
  makeCompoundType as _makeCompoundType,
  makeStructType as _makeStructType,
  makeType as _makeType,
  registerPackage as _registerPackage,
  stringType as _stringType,
  uint8Type as _uint8Type,
  valueType as _valueType,
} from '@attic/noms';
import type {
  Blob as _Blob,
  NomsSet as _NomsSet,
  Struct as _Struct,
  Value as _Value,
  float32 as _float32,
  float64 as _float64,
  uint8 as _uint8,
} from '@attic/noms';

const _pkg = new _Package([
  _makeStructType('StructWithUnionField',
      [
        new _Field('a', _float32Type, false),
      ],
      [
        new _Field('b', _float64Type, false),
        new _Field('c', _stringType, false),
        new _Field('d', _blobType, false),
        new _Field('e', _valueType, false),
        new _Field('f', _makeCompoundType(_Kind.Set, _uint8Type), false),
      ]
    ),
], [
]);
_registerPackage(_pkg);
export const typeForStructWithUnionField = _makeType(_pkg.ref, 0);
const StructWithUnionField$typeDef = _pkg.types[0];


type StructWithUnionField$Data = {
  a: _float32;
};

interface StructWithUnionField$Interface extends _Struct {
  constructor(data: StructWithUnionField$Data): void;
  a: _float32;  // readonly
  setA(value: _float32): StructWithUnionField$Interface;
  b: ?_float64;  // readonly
  setB(value: _float64): StructWithUnionField$Interface;
  c: ?string;  // readonly
  setC(value: string): StructWithUnionField$Interface;
  d: ?_Blob;  // readonly
  setD(value: _Blob): StructWithUnionField$Interface;
  e: ?_Value;  // readonly
  setE(value: _Value): StructWithUnionField$Interface;
  f: ?_NomsSet<_uint8>;  // readonly
  setF(value: _NomsSet<_uint8>): StructWithUnionField$Interface;
}

export const StructWithUnionField: Class<StructWithUnionField$Interface> = _createStructClass(typeForStructWithUnionField, StructWithUnionField$typeDef);
