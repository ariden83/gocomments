import tensorflow as tf
try:
    from tensorflow.python.compiler.tensorrt import trt_convert as trt
    print("TensorRT is available.")
except ImportError:
    print("TensorRT is not available.")
