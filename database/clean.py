import sqlite3
import os, glob

conn = sqlite3.connect("mindoc.db")
cur = conn.cursor()         #通过建立数据库游标对象，准备读写操作


cmd = """
SELECT 
    att.http_path
FROM
    md_attachment AS att
WHERE (att.document_id != 0 OR (NOT EXISTS( SELECT 1 FROM md_documents WHERE markdown LIKE ("%" || att.http_path || "%"))))
AND (att.document_id = 0 OR (NOT EXISTS( SELECT 1 FROM md_documents WHERE att.document_id = document_id )))
"""
cur.execute(cmd)
file_list = cur.fetchall()
for file_item in file_list:
    item_path = file_item[0]
    # 1. 删除os文件
    if os.path.exists(os.path.join("..", item_path[1:])):
        os.remove(os.path.join("..", item_path[1:]))
    
    # 2. 查询os是否删除成功，成功则删除附件记录
    if not os.path.exists(os.path.join("..", item_path[1:])):
        cmd = """
        delete
        from md_attachment
        WHERE http_path = '{}'
        """.format(item_path)
        cur.execute(cmd)
conn.commit()   #保存提交，确保数据保存成功

conn.close()        #关闭与数据库的连接