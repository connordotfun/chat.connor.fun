import metric_tree

class Filter:

    base_ban_list = ["apple", "orange", "flatiron"]
    base_allowed_list = ["the", "is", "a"]

    replacement_dict = { "@" : ("a", -.2), "!" : ("i", -.2), "&" : ("an", .2), "." : ("", -.05), "$" : ("s", -.2), "3" : ("e", -.2) }

    tree = metric_tree.Tree()
    tolerance = .85

    def __init__(self):
        for banned in self.base_ban_list:
            self.tree.insert_word(banned)

    def add_word(self, word):
        self.tree.insert_word(word)

    def censor_sentence(self, sentence, user_pcs=0):
        words = sentence.split(' ')
        new_pcs = user_pcs

        for word in words:
            modifier = 0
            check_word = word.lower()

            if word not in self.base_allowed_list:
                for i in range(len(check_word)):
                    if check_word[i] in self.replacement_dict:
                        ret_tup = self.replacement_dict[word[i]]
                        check_word = check_word[:i] + ret_tup[0] + check_word[i+1:]
                
                score = self.tree.get_score(check_word, modifier)

                if score > self.tolerance-(user_pcs/100):
                    sentence = sentence.replace(word, "*"*len(word))
                    new_pcs += 1
        
        return sentence, new_pcs


filt = Filter()

print(filt.censor_sentence("the @pple is orang3", 5)[0])