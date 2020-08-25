class InvalidParameterException(Exception):
  pass

class Form:
  pass

class Location:
  def __init__(self, *, x=None, y=None):
    if x is None:
      raise InvalidParameterException("x must be provided")
    if y is None:
      raise InvalidParameterException("y must be provided")

    self.X = x
    self.Y = y

class FontSize:
  def __init__(self, font_size):
    if font_size < 1:
      raise InvalidParameterException("font size must be a postive integer")
    self.font_size = font_size

class Size:
  def __init__(self, size):
    if size < 1:
      raise InvalidParameterException("size must be a postive integer")
    self.size = size