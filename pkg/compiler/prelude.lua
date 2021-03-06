-- prelude defines things that should
-- be available before any user code is run.

function _gi_GetRangeCheck(x, i)
  if x == nil or i < 0 or i >= #x then
     error("index out of range")
  end
  return x[i]
end;

function _gi_SetRangeCheck(x, i, val)
  --print("SetRangeCheck. x=", x, " i=", i, " val=", val)
  if x == nil or i < 0 or i >= #x then
     error("index out of range")
  end
  x[i] = val
  return val
end;

-- complex numbers
-- "complex[?]" ?

ffi = require('ffi')
point = ffi.metatype("struct { double re, im; }", {
    __add = function(a, b)
     return point(a.re + b.re, a.im + b.im)
 end
})

-- 1+2i
point = ffi.metatype("complex", {
    __add = function(a, b)
     return point(a.re + b.re, a.im + b.im)
 end
})

function _gi_NewComplex128(real, imag)

end
