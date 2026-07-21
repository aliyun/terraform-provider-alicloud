import logging  

def handler(event, context):
  logger = logging.getLogger()
  logger.info('hello world')
  return 'hello world'
