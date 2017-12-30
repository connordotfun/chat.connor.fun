import metric_tree
import re

class Filter:

    base_ban_list = ["apple", "orange", "flatiron"]
    base_allowed_list = ["the", "is", "a"]

    def __init__(self):
        tree = metric_tree.Tree()
        for banned in self.base_ban_list:
            tree.add_word(banned)

    def censor_sentence(self, input_sentence):
        words = input_sentence.split(' ', 1)

        
f = Filter()