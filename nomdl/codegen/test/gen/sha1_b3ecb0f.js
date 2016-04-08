// This file was generated by nomdl/codegen.
// @flow
/* eslint-disable */

import {
  Field as _Field,
  Kind as _Kind,
  Package as _Package,
  blobType as _blobType,
  createStructClass as _createStructClass,
  makeCompoundType as _makeCompoundType,
  makeStructType as _makeStructType,
  makeType as _makeType,
  registerPackage as _registerPackage,
} from '@attic/noms';
import type {
  Blob as _Blob,
  NomsList as _NomsList,
  Struct as _Struct,
} from '@attic/noms';

const _pkg = new _Package([
  _makeStructType('A',
      [
        new _Field('A', _makeCompoundType(_Kind.List, _makeCompoundType(_Kind.List, _blobType)), false),
      ],
      [

      ]
    ),
], [
]);
_registerPackage(_pkg);
export const typeForA = _makeType(_pkg.ref, 0);
const A$typeDef = _pkg.types[0];


type A$Data = {
  A: _NomsList<_NomsList<_Blob>>;
};

interface A$Interface extends _Struct {
  constructor(data: A$Data): void;
  A: _NomsList<_NomsList<_Blob>>;  // readonly
  setA(value: _NomsList<_NomsList<_Blob>>): A$Interface;
}

export const A: Class<A$Interface> = _createStructClass(typeForA, A$typeDef);
